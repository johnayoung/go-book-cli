package models

type LanguageModel interface {
	Generate(prompt string) (string, error)
	SetParameters(params map[string]interface{})
}
