package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ViolettaBykova/viot-tg-sirius/models/scenes"
	tb "gopkg.in/tucnak/telebot.v2"
)

var validIntervals = map[string]bool{
	"30 секунд": true,
	"1 минута":  true,
	"15 минут":  true,
	"1 час":     true,
	"6 часов":   true,
	"12 часов":  true,
}

type BotHandlers struct {
	botService BotService
	Bot        *tb.Bot
}

// NewBotHandlers создаёт новый экземпляр BotHandlers с зависимостями
func NewBotHandlers(bot *tb.Bot, botService BotService) *BotHandlers {
	return &BotHandlers{
		botService: botService,
		Bot:        bot,
	}
}

// HandleStart обрабатывает команду /start
func (h *BotHandlers) HandleStart(m *tb.Message) {
	ctx := context.TODO()
	err := h.botService.CreateUser(ctx, m.Sender.ID, "")
	if err != nil {
		h.Bot.Send(m.Sender, "Ошибка при сохранении пользователя.")
		return
	}

	h.Bot.Send(m.Sender, "Добро пожаловать! Пожалуйста, введите город для получения прогноза погоды.")
	h.botService.SetUserScene(ctx, m.Sender.ID, scenes.SceneEnterCity) // Устанавливаем сцену для ввода города
}

// HandleText обрабатывает текстовые сообщения
func (h *BotHandlers) HandleText(m *tb.Message) {
	ctx := context.TODO()
	scene, err := h.botService.GetUserScene(ctx, m.Sender.ID)
	if err != nil {
		h.Bot.Send(m.Sender, "Ошибка при получении состояния пользователя.")
		return
	}

	switch scene {
	case scenes.SceneEnterCity:
		h.botService.SetCity(ctx, m.Sender.ID, m.Text) // Сохраняем город
		h.Bot.Send(m.Sender, "Город сохранен. Выберите интервал обновления: 30 секунд, 1 минута, 15 минут, 1 час, 6 часов, 12 часов.")
		h.botService.SetUserScene(ctx, m.Sender.ID, scenes.SceneSelectInterval) // Переходим на сцену выбора интервала

	case scenes.SceneSelectInterval:

		// Проверка, что введенный интервал допустим
		if _, ok := validIntervals[m.Text]; !ok {
			h.Bot.Send(m.Sender, "Некорректный интервал. Пожалуйста, выберите один из следующих: 30 секунд, 1 минута, 15 минут, 1 час, 6 часов, 12 часов.")
			return // Прекращаем выполнение, если интервал некорректен
		}

		// Сохраняем интервал, если он корректен
		h.botService.SetUpdateInterval(ctx, m.Sender.ID, m.Text)
		if err := h.botService.ScheduleWeatherUpdate(ctx, m.Sender.ID, m.Text); err != nil {
			log.Printf("Ошибка планирования обновлений погоды для пользователя %d: %v", m.Sender.ID, err)
			h.Bot.Send(m.Sender, "Ошибка при планировании обновлений. Попробуйте еще раз.")
			return
		}
		h.Bot.Send(m.Sender, fmt.Sprintf("Интервал обновления установлен на %s.", m.Text))
		h.botService.SetUserScene(ctx, m.Sender.ID, scenes.SceneDefault)

	default:
		h.Bot.Send(m.Sender, "Команда не распознана. Используйте /start для начала.")
	}
}

// HandleSettings обрабатывает команду /settings
func (h *BotHandlers) HandleSettings(m *tb.Message) {
	ctx := context.TODO()
	h.Bot.Send(m.Sender, "Выберите интервал обновления: 30 секунд, 1 минута, 15 минут, 1 час, 6 часов, 12 часов.")
	h.botService.SetUserScene(ctx, m.Sender.ID, scenes.SceneSelectInterval)
}

// SendWeatherUpdate отправляет пользователю прогноз погоды
func (h *BotHandlers) SendWeatherUpdate(userID int64) {
	ctx := context.TODO()
	city, err := h.botService.GetUserCity(ctx, userID)
	if err != nil || city == "" {
		log.Printf("Не удалось найти город для пользователя %d\n", userID)
		return
	}

	weather := fetchWeatherData(city)
	h.Bot.Send(&tb.User{ID: userID}, fmt.Sprintf("Прогноз погоды для %s: %s", city, weather))
}

// fetchWeatherData – функция-заглушка для получения данных о погоде
func fetchWeatherData(city string) string {
	// Здесь может быть вызов API для получения данных о погоде
	return "ясно, +25°C"
}
