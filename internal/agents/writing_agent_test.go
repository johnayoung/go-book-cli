package agents

import (
	"testing"
)

type MockLanguageModel struct{}

func (m *MockLanguageModel) Generate(prompt string) (string, error) {
	return "mock response for: " + prompt, nil
}

func (m *MockLanguageModel) SetParameters(params map[string]interface{}) {}

func TestGenerateOutline(t *testing.T) {
	mockModel := &MockLanguageModel{}
	writingAgent := NewWritingAgent(mockModel)

	outline, err := writingAgent.GenerateOutline("Test Topic")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "mock response for: Generate a detailed outline for a book about Test Topic"
	if outline != expected {
		t.Errorf("Expected %v, got %v", expected, outline)
	}
}
