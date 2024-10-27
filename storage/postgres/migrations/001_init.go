package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	MigrationSet.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id uuid PRIMARY KEY,
            telegram_id BIGINT UNIQUE NOT NULL,
            city VARCHAR(100),
            update_interval VARCHAR(20) NOT NULL DEFAULT '1 час',
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			scene VARCHAR(50) NOT NULL DEFAULT 'default'
        );

`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
drop table users cascade;
;
`)
		return err
	})
}
