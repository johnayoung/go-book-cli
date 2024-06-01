package main

import (
	"go-book-ai/cmd/bookcli/cmd"
	"io"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("bookcli.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd.Execute()
}
