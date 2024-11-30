package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	Logger *slog.Logger
}

var instantiated *Logger = nil

func New() *Logger {
	if instantiated == nil {
		instantiated = &Logger{
			Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		}
	}
	return instantiated
}
