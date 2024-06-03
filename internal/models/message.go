package models

type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

type ConversationHistory struct {
	Messages []Message `json:"messages"`
}

func StartNewConversation(prompt string) ConversationHistory {
	return ConversationHistory{
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}
}

func (ch *ConversationHistory) AddMessage(role, content string) {
	ch.Messages = append(ch.Messages, Message{Role: role, Content: content})
}

func (ch *ConversationHistory) GetContext() []Message {
	return ch.Messages
}
