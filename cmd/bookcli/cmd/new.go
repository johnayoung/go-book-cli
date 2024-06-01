package cmd

import (
	"fmt"
	"log"
	"os"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/handlers"
	"go-book-ai/internal/models"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [topic]",
	Short: "Create a new book with the specified topic",
	Long:  `Create a new book and generate an outline based on the specified topic.`,
	Args:  cobra.ExactArgs(1),
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
		bookHandler := handlers.NewBookCommandHandler(writingAgent, reviewingAgent, fileManager, errorHandler)

		log.Printf("Starting to create a new book with topic: %s", topic)
		err := bookHandler.CreateNewBook(topic)
		if err != nil {
			log.Printf("Failed to create new book: %v", err)
			fmt.Fprintf(os.Stderr, "Failed to create new book: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
