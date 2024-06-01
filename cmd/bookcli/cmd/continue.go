package cmd

import (
	"fmt"
	"os"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/handlers"
	"go-book-ai/internal/models"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue [bookID]",
	Short: "Continue an existing book with the specified book ID",
	Long:  `Continue working on an existing book by providing the book ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bookID := args[0]

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

		err := bookHandler.ContinueExistingBook(bookID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to continue book: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(continueCmd)
}
