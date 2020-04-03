package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

/*
 * Custom layer on logrus logger
 */
var defaultLogger *Logger

// Logger is a... logger
type Logger struct {
	*logrus.Logger
}

// GetLogger returns default logger
func GetLogger() *Logger {
	if defaultLogger != nil {
		return defaultLogger
	}

	logger := &Logger{
		&logrus.Logger{},
	}

	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)

	envMinLevel := os.Getenv("LOGGER_MIN_LEVEL")

	if envMinLevel == "" {
		envMinLevel = "info"
	}

	minLevel, err := logrus.ParseLevel(envMinLevel)
	if err != nil {
		panic("invalid logger min level")
	}

	logger.SetLevel(minLevel)

	defaultLogger = logger

	return defaultLogger
}
