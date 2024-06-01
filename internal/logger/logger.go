package logger

import "log"

type Logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
}

type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Info(message string) {
	log.Println("INFO: " + message)
}

func (l *SimpleLogger) Warn(message string) {
	log.Println("WARN: " + message)
}

func (l *SimpleLogger) Error(message string) {
	log.Println("ERROR: " + message)
}
