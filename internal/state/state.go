package state

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type SubsectionState struct {
	Title     string `yaml:"title"`
	Generated bool   `yaml:"generated"`
}

type SectionState struct {
	Title       string            `yaml:"title"`
	Generated   bool              `yaml:"generated"`
	Subsections []SubsectionState `yaml:"subsections"`
}

type ChapterState struct {
	Title     string         `yaml:"title"`
	Generated bool           `yaml:"generated"`
	Sections  []SectionState `yaml:"sections"`
}

type State struct {
	OutlineGenerated bool           `yaml:"outlineGenerated"`
	Chapters         []ChapterState `yaml:"chapters"`
}

func NewState() *State {
	return &State{}
}

func (s *State) Save(path string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = os.WriteFile(filepath.Join(path, "state.yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save state file: %w", err)
	}

	return nil
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(filepath.Join(path, "state.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return NewState(), nil
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
