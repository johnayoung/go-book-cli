package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
)

// BookCommandHandler handles book-related commands.
type BookCommandHandler struct {
	WritingAgent   agents.WritingAgent
	ReviewingAgent agents.ReviewingAgent
	FileManager    *file.FileManager
	ErrorHandler   *errors.ErrorHandler
}

// NewBookCommandHandler returns a new BookCommandHandler.
func NewBookCommandHandler(writingAgent agents.WritingAgent, reviewingAgent agents.ReviewingAgent, fm *file.FileManager, eh *errors.ErrorHandler) *BookCommandHandler {
	return &BookCommandHandler{WritingAgent: writingAgent, ReviewingAgent: reviewingAgent, FileManager: fm, ErrorHandler: eh}
}

func (h *BookCommandHandler) CreateNewBook(topic string) error {
	bookPath := filepath.Join("books", topic)
	if _, err := os.Stat(bookPath); !os.IsNotExist(err) {
		return fmt.Errorf("a book with the topic '%s' already exists", topic)
	}

	err := os.MkdirAll(bookPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create book directory: %v", err)
	}

	outline, err := h.WritingAgent.GenerateOutline(topic)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to generate book outline: %v", err))
		if !h.ErrorHandler.HandleError(err) {
			return fmt.Errorf("retry attempts exhausted")
		}
	}

	outlinePath := filepath.Join(bookPath, "OUTLINE.md")
	err = h.FileManager.SaveOutline(outline, outlinePath)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to save outline file: %v", err))
		if !h.ErrorHandler.HandleError(err) {
			return fmt.Errorf("retry attempts exhausted")
		}
	}

	h.ErrorHandler.LogInfo(fmt.Sprintf("New book created with topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) ContinueExistingBook(bookID string) error {
	bookPath := filepath.Join("books", bookID)
	if _, err := os.Stat(bookPath); os.IsNotExist(err) {
		return fmt.Errorf("no book found with ID '%s'", bookID)
	}

	h.ErrorHandler.LogInfo(fmt.Sprintf("Continuing book with ID: %s", bookID))
	return nil
}
