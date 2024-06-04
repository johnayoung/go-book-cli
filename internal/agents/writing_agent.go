package agents

import (
	"fmt"
	"go-book-ai/internal/models"
	"go-book-ai/internal/outline"
)

type WritingAgent interface {
	GenerateOutline(topic string) (string, error)
	GenerateChapterOutline(chapterTitle string) (string, error)
	GenerateSectionContent(section outline.Section) (string, error)
	SendMessage(prompt string) (string, error)
}

type writingAgent struct {
	LanguageModel models.LanguageModel
}

func NewWritingAgent(model models.LanguageModel) WritingAgent {
	return &writingAgent{LanguageModel: model}
}

func (agent *writingAgent) GenerateOutline(topic string) (string, error) {
	prompt := fmt.Sprintf(`Generate a detailed book outline for a book titled "%s". 
The outline should include multiple chapters, each with several sections, formatted strictly in YAML as follows:

title: "%s"
chapters:
  - title: "Chapter 1: [Chapter Title]"
    sections:
      - title: "[Section Title]"
      - title: "[Section Title]"
  - title: "Chapter 2: [Chapter Title]"
    sections:
      - title: "[Section Title]"
      - title: "[Section Title]"
  ...
  - title: "Chapter N: [Chapter Title]"
    sections:
      - title: "[Section Title]"
      - title: "[Section Title]"

Please ensure the output is valid YAML and do not include any additional text or explanations.`, topic, topic)
	return prompt, nil
}

func (agent *writingAgent) GenerateChapterOutline(chapterTitle string) (string, error) {
	prompt := fmt.Sprintf(`Generate a detailed chapter outline for a chapter titled "%s". 
The outline should include sections and sub-sections, each with a brief description, formatted strictly in YAML as follows:

title: "%s"
sections:
  - title: "[Section Title]"
    description: "[Brief description of Section]"
    subsections:
      - title: "[Subsection Title]"
        description: "[Brief description of Subsection]"
      - title: "[Subsection Title]"
        description: "[Brief description of Subsection]"
  - title: "[Section Title]"
    description: "[Brief description of Section]"
  ...
  - title: "[Section Title]"
    description: "[Brief description of Section]"
    subsections:
      - title: "[Subsection Title]"
        description: "[Brief description of Subsection]"
      - title: "[Subsection Title]"
        description: "[Brief description of Subsection]"

Please ensure the output is valid YAML and do not include any additional text or explanations.`, chapterTitle, chapterTitle)
	return prompt, nil
}

func (agent *writingAgent) GenerateSectionContent(section outline.Section) (string, error) {
	subsectionsPrompt := ""
	for _, subsection := range section.Subsections {
		subsectionsPrompt += fmt.Sprintf("\n      - title: \"%s\"", subsection.Title)
	}

	prompt := fmt.Sprintf(`Write detailed content for the section titled "%s". 
Subsections:%s`, section.Title, subsectionsPrompt)

	return prompt, nil
}

func (agent *writingAgent) SendMessage(prompt string) (string, error) {
	agent.LanguageModel.SetParameters(map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	return agent.LanguageModel.Generate(prompt)
}
