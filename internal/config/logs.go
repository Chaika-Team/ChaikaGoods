package config

import (
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func ConfigureLogger(baseLogger log.Logger, levelStr string) log.Logger {
	var logger log.Logger
	_ = level.Info(baseLogger).Log("message", "Configuring logger", "level", levelStr)
	switch strings.ToLower(levelStr) {
	case "debug":
		logger = level.NewFilter(baseLogger, level.AllowDebug())
	case "info":
		logger = level.NewFilter(baseLogger, level.AllowInfo())
	case "warn":
		logger = level.NewFilter(baseLogger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(baseLogger, level.AllowError())
	default:
		// Default to info level if the provided level is invalid
		_ = level.Warn(baseLogger).Log("message", "Invalid log level, defaulting to info", "level", levelStr)
		logger = level.NewFilter(baseLogger, level.AllowInfo())
	}
	return logger
}
