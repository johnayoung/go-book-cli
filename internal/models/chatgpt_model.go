package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-book-ai/internal/errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ChatGPTModel struct {
	Parameters   map[string]interface{}
	ErrorHandler *errors.ErrorHandler
}

func NewChatGPTModel(errorHandler *errors.ErrorHandler) *ChatGPTModel {
	return &ChatGPTModel{ErrorHandler: errorHandler}
}

func (model *ChatGPTModel) SetParameters(params map[string]interface{}) error {
	model.Parameters = params
	return nil
}

func (model *ChatGPTModel) Generate(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set. Please set the OPENAI_API_KEY environment variable.")
	}

	url := "https://api.openai.com/v1/chat/completions"
	body := map[string]interface{}{
		"model":    "gpt-4",
		"messages": model.Parameters["messages"],
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}

	var resp *http.Response
	for retries := 0; retries <= model.ErrorHandler.RetryLimit; retries++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		model.ErrorHandler.LogError(fmt.Errorf("failed to execute request: %w", err))
		if !model.ErrorHandler.HandleError(err) {
			return "", fmt.Errorf("retry attempts exhausted")
		}
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("Request body: %s", jsonBody)
		return "", fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, string(responseBody))
	}

	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no choices in response body")
	}

	return respBody.Choices[0].Message.Content, nil
}
