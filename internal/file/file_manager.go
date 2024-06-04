package file

import (
	"fmt"
	"go-book-ai/internal/state"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type FileManager struct{}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) SaveSectionContent(content string, path string) error {
	log.Printf("DEBUG: Saving section content to %s", path) // Add debug logging
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Printf("DEBUG: Failed to save section content to %s: %v", path, err) // Add debug logging
		return fmt.Errorf("failed to save section content: %w", err)
	}
	log.Printf("DEBUG: Successfully saved section content to %s", path) // Add debug logging
	return nil
}

func (fm *FileManager) SaveState(path string, state interface{}) error {
	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}
	return nil
}

func (fm *FileManager) LoadState(path string) (*state.State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}
	var state state.State
	err = yaml.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return &state, nil
}
