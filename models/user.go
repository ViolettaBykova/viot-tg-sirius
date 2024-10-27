package models

import (
	"time"

	"github.com/ViolettaBykova/viot-tg-sirius/models/scenes"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel  `bun:"table:users"`
	ID             uuid.UUID    `bun:"id,pk,autoincrement"`
	TelegramID     int64        `bun:"telegram_id,unique,notnull"`
	City           string       `bun:"city"`
	UpdateInterval string       `bun:"update_interval,notnull,default:'1 час'"`
	CreatedAt      time.Time    `bun:"created_at,notnull,default:current_timestamp"`
	Scene          scenes.Scene `bun:"scene,notnull,default:'default'"` // Добавлено поле для состояния
}
