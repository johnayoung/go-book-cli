package file

import (
	"fmt"
	"go-book-ai/internal/logger"
	"go-book-ai/internal/state"
	"os"

	"gopkg.in/yaml.v2"
)

type FileManager struct {
	Logger logger.Logger
}

func NewFileManager(logger logger.Logger) *FileManager {
	return &FileManager{Logger: logger}
}

func (fm *FileManager) SaveSectionContent(content string, path string) error {
	fm.Logger.Debug(fmt.Sprintf("Saving section content to %s", path))
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fm.Logger.Debug(fmt.Sprintf("Failed to save section content to %s: %v", path, err))
		return fmt.Errorf("failed to save section content: %w", err)
	}
	fm.Logger.Debug(fmt.Sprintf("Successfully saved section content to %s", path))
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
		if os.IsNotExist(err) {
			// Create a new state if the file does not exist
			fm.Logger.Debug(fmt.Sprintf("State file %s does not exist, creating new state", path))
			newState := state.NewState()
			err := fm.SaveState(path, newState)
			if err != nil {
				return nil, fmt.Errorf("failed to create state file: %w", err)
			}
			return newState, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}
	var state state.State
	err = yaml.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return &state, nil
}
