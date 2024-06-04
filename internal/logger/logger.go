package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
}

type simpleLogger struct {
	debugEnabled bool
}

func NewSimpleLogger() Logger {
	debugEnabled := os.Getenv("DEBUG") == "true"
	return &simpleLogger{debugEnabled: debugEnabled}
}

func (l *simpleLogger) Info(message string) {
	log.Printf("INFO: %s", message)
}

func (l *simpleLogger) Debug(message string) {
	if l.debugEnabled {
		log.Printf("DEBUG: %s", message)
	}
}

func (l *simpleLogger) Error(message string) {
	log.Printf("ERROR: %s", message)
}
