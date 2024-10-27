package bot

import (
	"context"
	"log"

	"github.com/ViolettaBykova/viot-tg-sirius/models"
	"github.com/ViolettaBykova/viot-tg-sirius/models/scenes"
	"github.com/google/uuid"
)

// CreateUser создает нового пользователя или игнорирует, если он уже существует
func (s *Service) CreateUser(ctx context.Context, telegramID int64, city string) error {
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	id := uuid.New()
	user := &models.User{
		ID:         id,
		TelegramID: telegramID,
		City:       city,
	}

	err = s.store.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return err
	}

	log.Printf("User %d created successfully", telegramID)
	return s.store.TxCommit(ctx)
}

// GetUserCity возвращает город, указанный пользователем
func (s *Service) GetUserCity(ctx context.Context, telegramID int64) (string, error) {
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	city, err := s.store.GetUserCity(ctx, telegramID)
	if err != nil {
		log.Printf("Failed to get city for user %d: %v", telegramID, err)
		return "", err
	}

	return city, s.store.TxCommit(ctx)
}

func (s *Service) GetUserScene(ctx context.Context, telegramID int64) (scenes.Scene, error) {
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	user, err := s.store.GetUser(ctx, telegramID)
	if err != nil {
		return "", err
	}
	return user.Scene, s.store.TxCommit(ctx)
}

// SetUserScene устанавливает новую сцену для пользователя
func (s *Service) SetUserScene(ctx context.Context, telegramID int64, scene scenes.Scene) error {
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	err = s.store.UpdateUserScene(ctx, telegramID, scene)
	if err != nil {
		return err
	}
	return s.store.TxCommit(ctx)
}

func (s *Service) SetCity(ctx context.Context, telegramID int64, city string) error {
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	err = s.store.SetCity(ctx, telegramID, city)
	if err != nil {
		log.Printf("Failed to set city for user %d: %v", telegramID, err)
		return err
	}
	return s.store.TxCommit(ctx)
}

// SetUpdateInterval устанавливает интервал обновлений для пользователя
func (s *Service) SetUpdateInterval(ctx context.Context, telegramID int64, interval string) error {
	ctx, err := s.store.CtxWithTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = s.store.TxRollback(ctx)
	}()
	err = s.store.SetUpdateInterval(ctx, telegramID, interval)
	if err != nil {
		log.Printf("Failed to set update interval for user %d: %v", telegramID, err)
		return err
	}
	return s.store.TxCommit(ctx)
}
