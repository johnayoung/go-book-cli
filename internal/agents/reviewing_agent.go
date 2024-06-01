package agents

import "fmt"

type ReviewingAgent interface {
	ReviewContent(content string) (string, error)
}

type MockReviewingAgent struct{}

func NewMockReviewingAgent() *MockReviewingAgent {
	return &MockReviewingAgent{}
}

func (agent *MockReviewingAgent) ReviewContent(content string) (string, error) {
	// Simulate content review by appending a review note
	reviewedContent := fmt.Sprintf("%s\n\n[Reviewed content, with suggested improvements.]", content)
	return reviewedContent, nil
}
