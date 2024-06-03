package models

type MockLanguageModel struct{}

func (m *MockLanguageModel) Generate(prompt string) (string, error) {
	return "Generated content based on the prompt: " + prompt, nil
}

func (m *MockLanguageModel) SetParameters(params map[string]interface{}) {
	// Mock implementation does not need to set parameters
}
