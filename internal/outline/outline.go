package outline

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Subsection struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type Section struct {
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Subsections []Subsection `yaml:"subsections"`
}

type Chapter struct {
	Title    string    `yaml:"title"`
	Sections []Section `yaml:"sections"`
}

type Outline struct {
	Title    string    `yaml:"title"`
	Chapters []Chapter `yaml:"chapters"`
}

func NewOutline(title string) *Outline {
	return &Outline{Title: title}
}

func (o *Outline) Save(path string) error {
	data, err := yaml.Marshal(o)
	if err != nil {
		return fmt.Errorf("failed to marshal outline: %w", err)
	}

	err = os.WriteFile(filepath.Join(path, "OUTLINE.yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save outline file: %w", err)
	}

	return nil
}

func LoadOutline(path string) (*Outline, error) {
	data, err := os.ReadFile(filepath.Join(path, "OUTLINE.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read outline file: %w", err)
	}

	var outline Outline
	err = yaml.Unmarshal(data, &outline)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal outline: %w", err)
	}

	return &outline, nil
}

type ChapterOutline struct {
	Title    string    `yaml:"title"`
	Sections []Section `yaml:"sections"`
}

func NewChapterOutline(title string) *ChapterOutline {
	return &ChapterOutline{Title: title}
}

func (c *ChapterOutline) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal chapter outline: %w", err)
	}

	err = os.WriteFile(filepath.Join(path, "OUTLINE.yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save chapter outline file: %w", err)
	}

	return nil
}

func LoadChapterOutline(path string) (*ChapterOutline, error) {
	data, err := os.ReadFile(filepath.Join(path, "OUTLINE.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read chapter outline file: %w", err)
	}

	var chapterOutline ChapterOutline
	err = yaml.Unmarshal(data, &chapterOutline)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chapter outline: %w", err)
	}

	return &chapterOutline, nil
}
