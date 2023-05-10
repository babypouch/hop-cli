package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hop",
	Short: "Hop is a cli tool for executing babypouch workflows",
	Long:  `Hop is a cli tool for managing and executing baby pouch workflows built with love by kwngo and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
