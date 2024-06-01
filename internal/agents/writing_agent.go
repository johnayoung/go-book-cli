package agents

import (
	"fmt"
	"go-book-ai/internal/models"
)

type WritingAgent interface {
	GenerateOutline(topic string) (string, error)
	GenerateChapterOutline(chapterTitle string) (string, error)
	GenerateSectionContent(sectionTitle string) (string, error)
}

type writingAgent struct {
	LanguageModel models.LanguageModel
}

func NewWritingAgent(model models.LanguageModel) WritingAgent {
	return &writingAgent{LanguageModel: model}
}

func (agent *writingAgent) GenerateOutline(topic string) (string, error) {
	prompt := fmt.Sprintf(`Generate a detailed book outline for a book titled "%s". 
The outline should include multiple chapters, each with several sections, formatted in YAML as follows:

title: "%s"
chapters:
  - title: "Chapter 1"
    sections:
      - title: "Section 1.1"
      - title: "Section 1.2"
  - title: "Chapter 2"
    sections:
      - title: "Section 2.1"
      - title: "Section 2.2"
  ...
  - title: "Chapter N"
    sections:
      - title: "Section N.1"
      - title: "Section N.2"

Please ensure the output is valid YAML.`, topic, topic)
	return agent.LanguageModel.Generate(prompt)
}

func (agent *writingAgent) GenerateChapterOutline(chapterTitle string) (string, error) {
	prompt := fmt.Sprintf(`Generate a detailed chapter outline for a chapter titled "%s". 
The outline should include multiple sections, formatted in YAML as follows:

title: "%s"
sections:
  - title: "Section 1"
  - title: "Section 2"
  ...
  - title: "Section N"

Please ensure the output is valid YAML.`, chapterTitle, chapterTitle)
	return agent.LanguageModel.Generate(prompt)
}

func (agent *writingAgent) GenerateSectionContent(sectionTitle string) (string, error) {
	prompt := fmt.Sprintf(`Write detailed content for the section titled "%s".`, sectionTitle)
	return agent.LanguageModel.Generate(prompt)
}
