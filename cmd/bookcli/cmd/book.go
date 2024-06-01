package cmd

import (
	"fmt"
	"log"
	"os"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/handlers"
	"go-book-ai/internal/logger"
	"go-book-ai/internal/models"

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

		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "API key not set. Please set the OPENAI_API_KEY environment variable.\n")
			os.Exit(1)
		}

		chatGPTModel := models.NewChatGPTModel(apiKey)
		writingAgent := agents.NewWritingAgent(chatGPTModel)
		reviewingAgent := agents.NewMockReviewingAgent()
		fileManager := file.NewFileManager()
		errorHandler := errors.NewErrorHandler(3)
		logger := logger.NewSimpleLogger()
		bookHandler := handlers.NewBookCommandHandler(writingAgent, reviewingAgent, fileManager, errorHandler, logger)

		log.Printf("Starting process for book with topic: %s", topic)
		err := bookHandler.ProcessBook(topic)
		if err != nil {
			log.Printf("Failed to process book: %v", err)
			fmt.Fprintf(os.Stderr, "Failed to process book: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
}
