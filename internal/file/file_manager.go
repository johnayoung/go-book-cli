package file

import (
	"encoding/json"
	"fmt"
	"go-book-ai/internal/models"
	"os"
)

type FileManager struct{}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) SaveOutline(outline string, path string) error {
	return os.WriteFile(path, []byte(outline), 0644)
}

func (fm *FileManager) SaveSectionContent(content string, path string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to save section content: %w", err)
	}
	return nil
}

func (fm *FileManager) LoadOutline(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (fm *FileManager) LoadSectionContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (fm *FileManager) SaveHistoryToFile(history models.ConversationHistory, path string) error {
	data, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func (fm *FileManager) LoadHistoryFromFile(path string, history *models.ConversationHistory) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, history)
}

func (fm *FileManager) HandleCrashRecovery() error {
	// Placeholder for crash recovery logic
	// For example, check for incomplete files and handle them accordingly
	return nil
}
