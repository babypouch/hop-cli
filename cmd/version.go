package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hop",
	Long:  `All software has versions. This is Hop's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hop CLI V0.1 -- HEAD")
	},
}
