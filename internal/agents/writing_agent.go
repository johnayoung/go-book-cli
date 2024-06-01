package agents

import "go-book-ai/internal/models"

type WritingAgent struct {
	LanguageModel models.LanguageModel
}

func NewWritingAgent(model models.LanguageModel) *WritingAgent {
	return &WritingAgent{LanguageModel: model}
}

func (agent *WritingAgent) GenerateOutline(topic string) (string, error) {
	prompt := "Generate a detailed outline for a book about " + topic
	return agent.LanguageModel.Generate(prompt)
}

func (agent *WritingAgent) GenerateChapterOutline(chapterTitle string) (string, error) {
	prompt := "Generate a detailed outline for the chapter titled " + chapterTitle
	return agent.LanguageModel.Generate(prompt)
}

func (agent *WritingAgent) GenerateSectionContent(sectionTitle string) (string, error) {
	prompt := "Write detailed content for the section titled " + sectionTitle
	return agent.LanguageModel.Generate(prompt)
}
