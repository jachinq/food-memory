package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"food-memory/internal/config"
	"food-memory/internal/repository"
	"food-memory/internal/router"
	"food-memory/internal/storage"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	cfg := config.Load()
	cfg.UploadDir = absOr(cfg.UploadDir)
	cfg.PublicDir = firstExisting(cfg.PublicDir, "./public", "../frontend/dist")

	db, err := repository.Open(cfg)
	if err != nil {
		slog.Error("database", "err", err)
		os.Exit(1)
	}

	migDir := firstExisting(os.Getenv("MIGRATION_DIR"), "./migrations", "../backend/migrations")
	if err := repository.RunSQLMigrationsFromDir(db, migDir); err != nil {
		slog.Error("sql migrations", "err", err)
		os.Exit(1)
	}
	if err := repository.AutoMigrate(db); err != nil {
		slog.Error("auto migrate", "err", err)
		os.Exit(1)
	}

	fs, err := storage.NewLocal(cfg.UploadDir)
	if err != nil {
		slog.Error("storage", "err", err)
		os.Exit(1)
	}

	engine := router.New(cfg, db, fs)
	addr := ":" + strconv.Itoa(cfg.AppPort)
	slog.Info("食忆 API 已启动", "addr", addr, "upload", cfg.UploadDir)
	if err := engine.Run(addr); err != nil {
		slog.Error("server", "err", err)
		os.Exit(1)
	}
}

func absOr(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func firstExisting(candidates ...string) string {
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return ""
}
