package errors

import (
	"errors"
	"testing"
)

func TestHandleError(t *testing.T) {
	errorHandler := NewErrorHandler(3)

	// Simulate an error handling scenario
	retry := errorHandler.HandleError(errors.New("test error"))

	if retry {
		t.Errorf("Expected false, got %v", retry)
	}
}

func TestLogError(t *testing.T) {
	errorHandler := NewErrorHandler(3)

	// Simulate logging an error
	errorHandler.LogError(errors.New("test error"))
}

func TestLogInfo(t *testing.T) {
	errorHandler := NewErrorHandler(3)

	// Simulate logging an info message
	errorHandler.LogInfo("test info message")
}
