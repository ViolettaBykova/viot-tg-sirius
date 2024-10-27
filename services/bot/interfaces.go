package bot

import (
	"context"

	"github.com/ViolettaBykova/viot-tg-sirius/models"
	"github.com/ViolettaBykova/viot-tg-sirius/models/scenes"
)

type (
	Storage interface {
		// CtxWithTx returns a context with a transaction.
		// Creates a new transaction if there is no one and set in the context.
		CtxWithTx(ctx context.Context) (context.Context, error)
		// TxCommit commits the active transaction active from the context.
		// Returns nil if there is no transaction in the context.
		TxCommit(ctx context.Context) error
		// TxRollback rollbacks the active transaction active from the context.
		// Returns nil if there is no transaction in the context.
		TxRollback(ctx context.Context) error
		CreateUser(ctx context.Context, u *models.User) error
		GetUserCity(ctx context.Context, telegramID int64) (string, error)
		GetUser(ctx context.Context, telegramID int64) (*models.User, error)
		UpdateUserScene(ctx context.Context, telegramID int64, scene scenes.Scene) error
		SetCity(ctx context.Context, telegramID int64, city string) error
		SetUpdateInterval(ctx context.Context, telegramID int64, interval string) error
		GetAllUsersWithInterval(ctx context.Context) ([]models.User, error)
	}
)
