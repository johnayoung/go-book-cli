package cmd

import (
	"fmt"
	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/handlers"
	"go-book-ai/internal/logger"
	"go-book-ai/internal/models"
	"go-book-ai/internal/state"
	"go-book-ai/internal/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:   "book [topic]",
	Short: "Create or continue a book with the specified topic",
	Long: `Create a new book or continue an existing book by specifying the book topic.
If a book with the given topic already exists, the command will pick up where it left off.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]

		// Clean the topic name
		cleanedTopic := utils.CleanName(topic)

		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "API key not set. Please set the OPENAI_API_KEY environment variable.\n")
			os.Exit(1)
		}

		errorHandler := errors.NewErrorHandler(3)
		fileManager := file.NewFileManager()

		// Create the book directory if it does not exist
		bookPath := filepath.Join("books", cleanedTopic)
		err := os.MkdirAll(bookPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create book directory: %v", err)
		}

		stateFilePath := filepath.Join(bookPath, "state.yaml")

		// Load or create state
		bookState, err := fileManager.LoadState(stateFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				bookState = state.NewState()
				err = fileManager.SaveState(stateFilePath, bookState)
				if err != nil {
					log.Fatalf("Failed to create state file: %v", err)
				}
			} else {
				log.Fatalf("Failed to load state: %v", err)
			}
		}

		// Initialize agents and handlers
		chatGPTModel := models.NewChatGPTModel(errorHandler)
		writingAgent := agents.NewWritingAgent(chatGPTModel)
		reviewingAgent := agents.NewMockReviewingAgent()
		logger := logger.NewSimpleLogger()
		bookHandler := handlers.NewBookCommandHandler(writingAgent, reviewingAgent, fileManager, errorHandler, logger)

		log.Printf("Starting process for book with topic: %s", cleanedTopic)
		err = bookHandler.ProcessBook(topic)
		if err != nil {
			log.Printf("Failed to process book: %v", err)
			fmt.Fprintf(os.Stderr, "Failed to process book: %v\n", err)
			os.Exit(1)
		}

		// Save the updated state
		err = fileManager.SaveState(stateFilePath, bookState)
		if err != nil {
			log.Fatalf("Failed to save state: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
}
