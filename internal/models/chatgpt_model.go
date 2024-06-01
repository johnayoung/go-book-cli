package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ChatGPTModel struct {
	APIKey string
	Params map[string]interface{}
}

func NewChatGPTModel(apiKey string) *ChatGPTModel {
	return &ChatGPTModel{APIKey: apiKey, Params: make(map[string]interface{})}
}

func (model *ChatGPTModel) SetParameters(params map[string]interface{}) {
	model.Params = params
}

func (model *ChatGPTModel) Generate(prompt string) (string, error) {
	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
	requestBody, _ := json.Marshal(map[string]interface{}{
		"prompt":     prompt,
		"max_tokens": 1000,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", model.APIKey))

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to generate content: %s", resp.Status)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if text, ok := choices[0].(map[string]interface{})["text"].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("unexpected response format")
}
