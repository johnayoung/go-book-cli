package cmd

import (
	"fmt"
	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/handlers"
	"go-book-ai/internal/logger"
	"go-book-ai/internal/models"
	"go-book-ai/internal/utils"
	"os"

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
		logger := logger.NewSimpleLogger()
		if apiKey == "" {
			logger.Error("API key not set. Please set the OPENAI_API_KEY environment variable.")
			os.Exit(1)
		}

		errorHandler := errors.NewErrorHandler(3)
		fileManager := file.NewFileManager(logger)

		chatGPTModel := models.NewChatGPTModel(errorHandler)
		writingAgent := agents.NewWritingAgent(chatGPTModel)
		reviewingAgent := agents.NewMockReviewingAgent()
		bookHandler := handlers.NewBookCommandHandler(writingAgent, reviewingAgent, fileManager, errorHandler, logger)

		logger.Info(fmt.Sprintf("Starting process for book with topic: %s", cleanedTopic))
		err := bookHandler.ProcessBook(cleanedTopic)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to process book: %v", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
}
