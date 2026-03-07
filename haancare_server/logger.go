package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger() {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
}
