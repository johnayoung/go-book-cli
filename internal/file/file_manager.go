package file

import (
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
	return os.WriteFile(path, []byte(content), 0644)
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

func (fm *FileManager) HandleCrashRecovery() error {
	// Placeholder for crash recovery logic
	// For example, check for incomplete files and handle them accordingly
	return nil
}
