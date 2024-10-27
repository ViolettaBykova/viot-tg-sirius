package handlers

import (
	"context"

	"github.com/ViolettaBykova/viot-tg-sirius/models/scenes"
)

type (
	BotService interface {
		CreateUser(ctx context.Context, telegramID int64, city string) error
		GetUserCity(ctx context.Context, telegramID int64) (string, error)
		GetUserScene(ctx context.Context, telegramID int64) (scenes.Scene, error)
		SetUserScene(ctx context.Context, telegramID int64, scene scenes.Scene) error
		SetCity(ctx context.Context, telegramID int64, city string) error
		SetUpdateInterval(ctx context.Context, telegramID int64, interval string) error

		ScheduleWeatherUpdate(ctx context.Context, telegramID int64, interval string) error
	}
)
