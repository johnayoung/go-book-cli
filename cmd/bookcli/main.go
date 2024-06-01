package main

import (
	"fmt"
	"os"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/handlers"
	"go-book-ai/internal/models"
)

func main() {
	// Get API key from environment variable or configuration
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "API key not set. Please set the OPENAI_API_KEY environment variable.\n")
		os.Exit(1)
	}

	// Initialize the language model
	chatGPTModel := models.NewChatGPTModel(apiKey)

	// Initialize the agents
	writingAgent := agents.NewWritingAgent(chatGPTModel)
	reviewingAgent := agents.NewMockReviewingAgent()
	fileManager := file.NewFileManager()
	errorHandler := errors.NewErrorHandler(3) // Retry limit set to 3
	bookHandler = handlers.NewBookCommandHandler(writingAgent, reviewingAgent, fileManager, errorHandler)

	// Initialize the application
	err := initializeApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Parse command-line arguments
	err = parseArguments(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %v\n", err)
		os.Exit(1)
	}
}
