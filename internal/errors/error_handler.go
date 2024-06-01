package errors

import (
	"log"
	"time"
)

type ErrorHandler struct {
	RetryLimit int
}

func NewErrorHandler(retryLimit int) *ErrorHandler {
	return &ErrorHandler{RetryLimit: retryLimit}
}

func (eh *ErrorHandler) HandleError(err error) bool {
	log.Printf("Error: %v", err)

	for retries := 0; retries < eh.RetryLimit; retries++ {
		log.Printf("Retrying... Attempt %d/%d", retries+1, eh.RetryLimit)
		time.Sleep(2 * time.Second) // Backoff before retrying

		// Here, you could add code to reattempt the operation that caused the error
		// Return true if retry was successful
	}

	return false // Return false if all retries fail
}

func (eh *ErrorHandler) LogError(err error) {
	log.Printf("Error: %v", err)
}

func (eh *ErrorHandler) LogInfo(message string) {
	log.Println(message)
}
