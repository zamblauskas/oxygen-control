package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Setup() zerolog.Logger {
	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	fileWriter := &lumberjack.Logger{
		Filename:   "logs/oxygen-control.log",
		MaxSize:    10, // megabytes
		MaxBackups: 7,
		MaxAge:     1, // days
		Compress:   true,
	}

	// Combine both writers
	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	// Set global logger
	log.Logger = logger

	return logger
}
