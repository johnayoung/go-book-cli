package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BookDir   string `yaml:"book_dir"`
	OpenAIKey string `yaml:"openai_key"`
	LogLevel  string `yaml:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		config.OpenAIKey = apiKey
	}

	return &config, nil
}
