// Package database owns infrastructure connections.
package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/config"
)

const pingTimeout = 5 * time.Second

// ConnectPostgres opens and verifies a PostgreSQL connection using the
// environment-backed application config.
func ConnectPostgres(cfg config.Config) (*gorm.DB, error) {
	dsn := buildPostgresDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	if err := PingPostgres(context.Background(), db); err != nil {
		return nil, fmt.Errorf("verify PostgreSQL connection: %w", err)
	}

	return db, nil
}

func buildPostgresDSN(cfg config.Config) string {
	if databaseURL := strings.TrimSpace(cfg.DatabaseURL); databaseURL != "" {
		return databaseURL
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}

// PingPostgres verifies that GORM's underlying SQL connection can reach
// PostgreSQL. It is shared by startup and the database health endpoint.
func PingPostgres(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get SQL database handle: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return nil
}
