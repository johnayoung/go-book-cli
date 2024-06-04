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
		subsectionsPrompt += fmt.Sprintf("\n- title: \"%s\"\n  description: \"[Detailed description of the subsection]\"", subsection.Title)
	}

	prompt := fmt.Sprintf(`You are writing a detailed section for a book. The section is titled "%s" and it contains the following subsections:
%s

Please write a comprehensive draft for this section in Markdown format. The content should include:

1. An introduction that provides an overview of the section.
2. Detailed explanations for each of the subsections listed, with clear and thorough descriptions.
3. Practical examples or case studies where relevant.
4. Conclusion that summarizes the key points covered in the section.

Make sure the content is engaging, informative, and suitable for a book. Write in a clear and professional tone, and ensure the output is well-structured and coherent. Use markdown formatting including headings, subheadings, lists, code blocks, and other formatting features where appropriate.`, section.Title, subsectionsPrompt)

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
