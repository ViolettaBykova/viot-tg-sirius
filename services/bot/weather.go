package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

// StartScheduler запускает cron планировщик
func (s *Service) StartScheduler() {
	s.cron.Start()
}

// StopScheduler останавливает cron планировщик
func (s *Service) StopScheduler() {
	s.cron.Stop()
}

func (s *Service) loadScheduledJobs() {
	ctx := context.Background()

	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = s.store.TxRollback(ctx)
		}

	}()
	users, err := s.store.GetAllUsersWithInterval(ctx)
	if err != nil {
		log.Printf("Ошибка при загрузке пользователей для cron задач: %v", err)
		return
	}
	if err := s.store.TxCommit(ctx); err != nil {
		return
	}

	for _, user := range users {
		if user.UpdateInterval != "" {
			// Планируем задачу cron для каждого пользователя
			if err := s.ScheduleWeatherUpdate(ctx, user.TelegramID, user.UpdateInterval); err != nil {
				log.Printf("Ошибка планирования обновлений для пользователя %d: %v", user.TelegramID, err)
			}
		}
	}

}

// ScheduleWeatherUpdate создает или обновляет задачу для пользователя
func (s *Service) ScheduleWeatherUpdate(ctx context.Context, telegramID int64, interval string) error {
	// Удаляем существующую задачу для пользователя, если она есть
	if entryID, exists := s.cronJobs[telegramID]; exists {
		s.cron.Remove(cron.EntryID(entryID))
	}

	// Получаем спецификацию cron для интервала
	cronSpec, err := getCronSpec(interval)
	if err != nil {
		return err
	}

	// Создаем новую задачу для обновления погоды
	entryID, err := s.cron.AddFunc(cronSpec, func() {
		err := s.sendWeatherUpdate(ctx, telegramID)
		if err != nil {
			log.Printf("Error sending weather update for user %d: %v", telegramID, err)
			s.bot.Send(&tb.User{ID: telegramID}, "Ошибка получение данных")
		}
	})
	if err != nil {
		return err
	}

	// Сохраняем ID задачи
	s.cronJobs[telegramID] = int(entryID)
	return nil
}

// sendWeatherUpdate отправляет сообщение с прогнозом погоды пользователю
func (s *Service) sendWeatherUpdate(ctx context.Context, telegramID int64) error {
	ctx = context.WithValue(ctx, "tx", nil)
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	city, err := s.store.GetUserCity(ctx, telegramID)
	if err != nil || city == "" {
		return fmt.Errorf("city not found for user %d: %v", telegramID, err)
	}
	if err := s.store.TxCommit(ctx); err != nil {
		return err
	}
	// Получаем данные о погоде с помощью weatherAPI
	weatherData, err := s.weatherAPI.GetCurrentWeather(city)
	if err != nil {
		log.Printf("Ошибка при получении данных о погоде: %v", err)
		return err
	}

	// Формируем сообщение и отправляем его пользователю
	message := fmt.Sprintf("Погода в %s: %s\nТемпература: %.1f°C\nВлажность: %d%%", weatherData.Name, weatherData.Weather[0].Description, weatherData.Main.Temp, weatherData.Main.Humidity)
	_, err = s.bot.Send(&tb.User{ID: telegramID}, message)

	return err

}

// getCronSpec возвращает выражение cron для заданного интервала
func getCronSpec(interval string) (string, error) {
	switch interval {
	case "30 секунд":
		return "@every 30s", nil
	case "1 минута":
		return "@every 1m", nil
	case "15 минут":
		return "@every 15m", nil
	case "1 час":
		return "@every 1h", nil
	case "6 часов":
		return "@every 6h", nil
	case "12 часов":
		return "@every 12h", nil
	default:
		return "@every 1h", nil
	}
}

// fetchWeatherData – функция-заглушка для получения данных о погоде
func fetchWeatherData(city string) string {
	// Здесь будет вызов API для получения данных о погоде
	return "ясно, +25°C"
}
