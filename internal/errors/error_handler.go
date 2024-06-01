package errors

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
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
		waitTime := eh.exponentialBackoff(retries)
		log.Printf("Retrying... Attempt %d/%d after %v", retries+1, eh.RetryLimit, waitTime)
		time.Sleep(waitTime)

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

func (eh *ErrorHandler) exponentialBackoff(attempt int) time.Duration {
	min := 100  // 100ms
	max := 1000 // 1000ms
	backoff := min * int(math.Pow(2, float64(attempt)))
	if backoff > max {
		backoff = max
	}
	return time.Duration(backoff+rand.Intn(min)) * time.Millisecond
}

func (eh *ErrorHandler) IsNetworkError(err error) bool {
	if err, ok := err.(*url.Error); ok {
		return true
	}
	return false
}

func (eh *ErrorHandler) IsRateLimitError(resp *http.Response) bool {
	return resp.StatusCode == http.StatusTooManyRequests
}
