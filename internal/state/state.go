package state

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type SubsectionState struct {
	Title string `yaml:"title"`
}

type SectionState struct {
	Title          string            `yaml:"title"`
	DraftGenerated bool              `yaml:"draft_generated"`
	Subsections    []SubsectionState `yaml:"subsections"`
}

type ChapterState struct {
	Title            string         `yaml:"title"`
	OutlineGenerated bool           `yaml:"outline_generated"`
	DraftGenerated   bool           `yaml:"draft_generated"`
	Sections         []SectionState `yaml:"sections"`
}

type State struct {
	OutlineGenerated bool           `yaml:"outline_generated"`
	Chapters         []ChapterState `yaml:"chapters"`
	MessageHistory   []Message      `yaml:"message_history"`
}

type Message struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

func NewState() *State {
	return &State{
		OutlineGenerated: false,
		Chapters:         []ChapterState{},
		MessageHistory:   []Message{},
	}
}

func (s *State) Save(path string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}
	return nil
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new state if the file does not exist
			newState := NewState()
			err := newState.Save(path)
			if err != nil {
				return nil, fmt.Errorf("failed to create state file: %w", err)
			}
			return newState, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}
	var state State
	err = yaml.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return &state, nil
}
