package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ViolettaBykova/viot-tg-sirius/storage/postgres/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
)

type Storage struct {
	db *bun.DB
}

func New(address, database, user, password string, pool int, envName string) (*Storage, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(address),
		pgdriver.WithInsecure(true),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database),
		pgdriver.WithApplicationName("gift-cards"),
		pgdriver.WithTimeout(30*time.Second),
	))

	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	db.SetMaxOpenConns(pool)

	if envName == "dev" || envName == "stage" || envName == "testing" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil

}

// CheckHealth pings database to make sure it's up and running.
func (s *Storage) CheckHealth(ctx context.Context) error {
	if s == nil {
		return fmt.Errorf("postgres is not initialized")
	}
	err := s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("connection to the database is failed: %w", err)
	}

	return nil
}

// Migrate checks if db initialized and applies migrations.
func (s *Storage) Migrate() error {
	if s.db == nil {
		return errors.New("database is not initialized")
	}
	if migrator, err := runInitCommand(s.db); err != nil {
		return err
	} else {
		_, err := migrator.Migrate(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
}

// runInitCommand initializes new migrator.
func runInitCommand(db *bun.DB) (*migrate.Migrator, error) {
	migrator := migrate.NewMigrator(db, migrations.MigrationSet)
	if err := migrator.Init(context.Background()); err != nil {
		return nil, err
	}
	return migrator, nil
}

// Close shuts down current connection.
func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

// DB method used in tests to get sql.DB connection.
func (s *Storage) DB() (*sql.DB, error) {
	return s.db.DB, nil
}
