package agents

import "go-book-ai/internal/models"

// WritingAgent defines the methods for generating book content.
type WritingAgent interface {
	GenerateOutline(topic string) (string, error)
	GenerateChapterOutline(chapterTitle string) (string, error)
	GenerateSectionContent(sectionTitle string) (string, error)
}

// writingAgent implements the WritingAgent interface.
type writingAgent struct {
	LanguageModel models.LanguageModel
}

// NewWritingAgent returns an instance of writingAgent.
func NewWritingAgent(model models.LanguageModel) WritingAgent {
	return &writingAgent{LanguageModel: model}
}

func (agent *writingAgent) GenerateOutline(topic string) (string, error) {
	prompt := "Generate a detailed outline for a book about " + topic
	return agent.LanguageModel.Generate(prompt)
}

func (agent *writingAgent) GenerateChapterOutline(chapterTitle string) (string, error) {
	prompt := "Generate a detailed outline for the chapter titled " + chapterTitle
	return agent.LanguageModel.Generate(prompt)
}

func (agent *writingAgent) GenerateSectionContent(sectionTitle string) (string, error) {
	prompt := "Write detailed content for the section titled " + sectionTitle
	return agent.LanguageModel.Generate(prompt)
}
