package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

func Init(settings *viper.Viper) *slog.Logger {
	lvl := slog.LevelInfo
	level := strings.ToLower(settings.GetString("log_level"))
	env := strings.ToLower(settings.GetString("env"))
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	var logger *slog.Logger
	if env == "production" {
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: lvl == slog.LevelError,
			Level:     lvl,
		})
		logger = slog.New(handler)
	} else {
		opts := tint.Options{
			AddSource: false,
			Level:     lvl,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Value.Kind() == slog.KindAny {
					if _, ok := a.Value.Any().(error); ok {
						return tint.Attr(9, a)
					}
				}
				return a
			},
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		logger = slog.New(
			tint.NewHandler(os.Stdout, &opts),
		)

	}
	slog.SetDefault(logger)
	return logger
}
