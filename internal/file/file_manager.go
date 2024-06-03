package file

import (
	"fmt"
	"go-book-ai/internal/state"
	"os"

	"gopkg.in/yaml.v2"
)

type FileManager struct{}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) SaveState(path string, bookState *state.State) error {
	data, err := yaml.Marshal(bookState)
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
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var bookState state.State
	err = yaml.Unmarshal(data, &bookState)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &bookState, nil
}

func (fm *FileManager) SaveSectionContent(content string, path string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write section content: %w", err)
	}
	return nil
}
