package repository

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"food-memory/internal/config"
	"food-memory/internal/model"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Dish{},
		&model.CookRecord{},
		&model.Tag{},
		&model.DishTag{},
		&model.Attachment{},
		&model.RecookPlan{},
		&model.AppSetting{},
	)
}

func RunSQLMigrationsFromDir(db *gorm.DB, dir string) error {
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("migration dir not found, skip sql migrations", "dir", dir)
			return nil
		}
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		var count int64
		if err := db.Raw("SELECT COUNT(1) FROM schema_migrations WHERE filename = ?", name).Scan(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if err := db.Exec(string(body)).Error; err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if err := db.Exec("INSERT INTO schema_migrations (filename) VALUES (?)", name).Error; err != nil {
			return err
		}
		slog.Info("applied migration", "file", name)
	}
	return nil
}
