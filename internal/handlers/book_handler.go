package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/outline"
	"go-book-ai/internal/state"
	"go-book-ai/internal/utils"

	"gopkg.in/yaml.v2"
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
	folderName := utils.CleanName(topic)
	bookPath := filepath.Join("books", folderName)
	log.Printf("Creating new book folder: %s", bookPath)

	if _, err := os.Stat(bookPath); !os.IsNotExist(err) {
		state, err := state.LoadState(bookPath)
		if err != nil {
			return fmt.Errorf("failed to load state: %v", err)
		}

		if state.OutlineGenerated {
			log.Printf("Outline already generated for book: %s", topic)
			return nil
		}

		return fmt.Errorf("a book with the topic '%s' already exists", topic)
	}

	err := os.MkdirAll(bookPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create book directory: %v", err)
	}

	state := state.NewState()
	bookOutline := outline.NewOutline(topic)

	log.Println("Generating book outline...")
	outlineContent, err := h.WritingAgent.GenerateOutline(topic)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to generate book outline: %v", err))
		if !h.ErrorHandler.HandleError(err) {
			return fmt.Errorf("retry attempts exhausted")
		}
	}

	// Log the generated outline content to understand its format
	log.Printf("Generated outline content:\n%s", outlineContent)

	// Assuming the generated outline content is in YAML format,
	// we need to parse it and populate the `Outline` struct.
	err = yaml.Unmarshal([]byte(outlineContent), bookOutline)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to parse generated outline: %v", err))
		return fmt.Errorf("failed to parse generated outline: %v", err)
	}

	err = bookOutline.Save(bookPath)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to save outline: %v", err))
		return fmt.Errorf("failed to save outline: %v", err)
	}

	state.OutlineGenerated = true
	err = state.Save(bookPath)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to save state: %v", err))
		return fmt.Errorf("failed to save state: %v", err)
	}

	h.ErrorHandler.LogInfo(fmt.Sprintf("New book created with topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) ContinueExistingBook(bookID string) error {
	bookPath := filepath.Join("books", bookID)
	log.Printf("Continuing book with ID: %s at path: %s", bookID, bookPath)

	if _, err := os.Stat(bookPath); os.IsNotExist(err) {
		return fmt.Errorf("no book found with ID '%s'", bookID)
	}

	state, err := state.LoadState(bookPath)
	if err != nil {
		return fmt.Errorf("failed to load state: %v", err)
	}

	if !state.OutlineGenerated {
		log.Println("Generating book outline...")
		outlineContent, err := h.WritingAgent.GenerateOutline(bookID)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to generate book outline: %v", err))
			if !h.ErrorHandler.HandleError(err) {
				return fmt.Errorf("retry attempts exhausted")
			}
		}

		// Log the generated outline content to understand its format
		log.Printf("Generated outline content:\n%s", outlineContent)

		// Assuming the generated outline content is in YAML format,
		// we need to parse it and populate the `Outline` struct.
		bookOutline := outline.NewOutline(bookID)
		err = yaml.Unmarshal([]byte(outlineContent), bookOutline)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to parse generated outline: %v", err))
			return fmt.Errorf("failed to parse generated outline: %v", err)
		}

		err = bookOutline.Save(bookPath)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to save outline: %v", err))
			return fmt.Errorf("failed to save outline: %v", err)
		}

		state.OutlineGenerated = true
		err = state.Save(bookPath)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to save state: %v", err))
			return fmt.Errorf("failed to save state: %v", err)
		}
	}

	h.ErrorHandler.LogInfo(fmt.Sprintf("Continuing book with ID: %s", bookID))
	return nil
}
