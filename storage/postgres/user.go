package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ViolettaBykova/viot-tg-sirius/models"
	"github.com/ViolettaBykova/viot-tg-sirius/models/scenes"
)

func (s *Storage) CreateUser(ctx context.Context, u *models.User) error {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return ErrTxNotFound
	}

	if _, err := tx.NewInsert().Model(u).On("CONFLICT (telegram_id) DO NOTHING").Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserCity(ctx context.Context, telegramID int64) (string, error) {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return "", ErrTxNotFound
	}

	var user models.User
	err := tx.NewSelect().
		Model(&user).Where("telegram_id = ?", telegramID).Column("city").Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "nil", nil
		}
		log.Printf("Error fetching user city: %v\n", err)
		return "", err
	}
	return user.City, nil
}

func (s *Storage) UpdateUserScene(ctx context.Context, telegramID int64, scene scenes.Scene) error {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return ErrTxNotFound
	}
	_, err := tx.NewUpdate().
		Model(&models.User{}).
		Set("scene = ?", scene).
		Where("telegram_id = ?", telegramID).
		Exec(context.Background())
	return err
}

// GetUser возвращает пользователя по telegram_id
func (s *Storage) GetUser(ctx context.Context, telegramID int64) (*models.User, error) {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return nil, ErrTxNotFound
	}
	var user models.User
	err := tx.NewSelect().Model(&user).Where("telegram_id = ?", telegramID).Scan(context.Background())
	if err != nil {
		log.Printf("Error fetching user: %v\n", err)
		return nil, err
	}
	return &user, nil
}

// SetCity устанавливает город для пользователя
func (s *Storage) SetCity(ctx context.Context, telegramID int64, city string) error {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return ErrTxNotFound
	}

	_, err := tx.NewUpdate().
		Model(&models.User{}).
		Set("city = ?", city).
		Where("telegram_id = ?", telegramID).
		Exec(ctx)
	return err
}

// SetUpdateInterval устанавливает интервал обновлений для пользователя
func (s *Storage) SetUpdateInterval(ctx context.Context, telegramID int64, interval string) error {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return ErrTxNotFound
	}

	_, err := tx.NewUpdate().
		Model(&models.User{}).
		Set("update_interval = ?", interval).
		Where("telegram_id = ?", telegramID).
		Exec(ctx)
	return err
}

func (s *Storage) GetAllUsersWithInterval(ctx context.Context) ([]models.User, error) {
	tx, ok := txFromCtx(ctx)
	if !ok {
		return nil, ErrTxNotFound
	}
	var users []models.User
	err := tx.NewSelect().
		Model(&users).
		Column("telegram_id", "update_interval").
		Where("update_interval IS NOT NULL AND update_interval != ''").
		Scan(ctx)
	if err != nil {
		log.Printf("Ошибка при получении пользователей с интервалом: %v", err)
		return nil, err
	}
	return users, nil
}
