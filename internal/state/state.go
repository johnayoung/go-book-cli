package state

type Message struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

type SubsectionState struct {
	Title string `yaml:"title"`
}

type SectionState struct {
	Title            string            `yaml:"title"`
	OutlineGenerated bool              `yaml:"outline_generated"`
	DraftGenerated   bool              `yaml:"draft_generated"`
	Subsections      []SubsectionState `yaml:"subsections"`
}

type ChapterState struct {
	Title            string         `yaml:"title"`
	OutlineGenerated bool           `yaml:"outline_generated"`
	DraftGenerated   bool           `yaml:"draft_generated"`
	Sections         []SectionState `yaml:"sections"`
}

type State struct {
	OutlineGenerated bool           `yaml:"outline_generated"`
	Chapters         []ChapterState `yaml:"chapters"`
	MessageHistory   []Message      `yaml:"message_history"`
}

func NewState() *State {
	return &State{
		OutlineGenerated: false,
		Chapters:         []ChapterState{},
		MessageHistory:   []Message{},
	}
}

func GetContext(history []Message) []map[string]string {
	context := make([]map[string]string, len(history))
	for i, msg := range history {
		context[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}
	return context
}
