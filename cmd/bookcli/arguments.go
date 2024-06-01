package main

import (
	"errors"
	"fmt"
	"go-book-ai/internal/agents"
	"go-book-ai/internal/handlers"
)

var writingAgent = agents.NewMockWritingAgent()
var bookHandler = handlers.NewBookCommandHandler(writingAgent)

func parseArguments(args []string) error {
	if len(args) < 2 {
		return errors.New("not enough arguments")
	}

	command := args[1]

	switch command {
	case "new":
		if len(args) < 3 {
			return errors.New("missing topic for new book")
		}
		topic := args[2]
		return handleCreateNewBook(topic)
	case "continue":
		if len(args) < 3 {
			return errors.New("missing book ID to continue")
		}
		bookID := args[2]
		return handleContinueExistingBook(bookID)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleCreateNewBook(topic string) error {
	return bookHandler.CreateNewBook(topic)
}

func handleContinueExistingBook(bookID string) error {
	return bookHandler.ContinueExistingBook(bookID)
}
