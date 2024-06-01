package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookcli",
	Short: "bookcli is a CLI tool for generating book chapters using AI.",
	Long:  `bookcli leverages AI to create detailed outlines and content for book chapters.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
