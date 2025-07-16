package config

import (
	"log/slog"

	slogbetterstack "github.com/samber/slog-betterstack"
)

func InitLogger(c *Config) {
	var logger *slog.Logger
	var logLevel slog.Level

	switch c.Log.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	if c.Log.Token != "" && c.Log.Endpoint != "" {
		logger = slog.New(slogbetterstack.Option{
			Token:    c.Log.Token,
			Endpoint: c.Log.Endpoint,
			Level:    logLevel,
		}.NewBetterstackHandler())
	} else {
		logger = slog.Default()
	}

	slog.SetLogLoggerLevel(logLevel)
	slog.SetDefault(logger)
}
