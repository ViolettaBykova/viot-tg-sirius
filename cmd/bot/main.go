package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ViolettaBykova/viot-tg-sirius/handlers"
	"github.com/ViolettaBykova/viot-tg-sirius/pkg/weather"
	botservice "github.com/ViolettaBykova/viot-tg-sirius/services/bot"
	"github.com/ViolettaBykova/viot-tg-sirius/storage/postgres"
	"github.com/spf13/viper"

	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	// Загружаем переменные окружения из .env
	viper.SetConfigFile(".env") // Устанавливаем файл конфигурации
	viper.AutomaticEnv()        // Автоматически считываем переменные из среды

	// Чтение конфигурации из файла .env
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Не удалось прочитать файл .env: %v\n", err)
	}
	// STORAGE
	fmt.Println(1, os.Getenv("DATABASE_ADDR"), 1)
	db, err := postgres.New(
		viper.GetString("DATABASE_ADDR"),
		viper.GetString("DATABASE_NAME"),
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		10,
		"dev",
	)
	if err != nil {
		log.Fatal(err, "init db")
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err, "migrating db")
	}

	// Initialize Telegram bot
	bot, err := tb.NewBot(tb.Settings{
		Token:  viper.GetString("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Failed to create bot: %v\n", err)
	}
	cronScheduler := cron.New()
	apiKey := viper.GetString("OPENWEATHER_API_KEY")
	weatherClient := weather.New(apiKey)

	// Создание botService с weatherClient
	botService := botservice.New(db, bot, cronScheduler, weatherClient) // Создаем botService с cron
	botService.StartScheduler()

	botHandlers := handlers.NewBotHandlers(bot, botService)

	// Set up bot handlers
	bot.Handle("/start", botHandlers.HandleStart)
	bot.Handle("/settings", botHandlers.HandleSettings)
	bot.Handle(tb.OnText, botHandlers.HandleText)

	// Start the bot
	log.Println("Бот запущен...")
	bot.Start()
}
