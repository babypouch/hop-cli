package cmd

import (
	"fmt"
	"os"

	"github.com/babypouch/hop-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "hop",
	Short: "Hop is a cli tool for executing babypouch workflows",
	Long:  `Hop is a cli tool for managing and executing baby pouch workflows built with love by kwngo and friends in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize(initConfig)

}

type HopConfig struct {
	AuthToken       string
	Host            string
	MicrolinkApiKey string
	SerpAuthToken   string
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	cfgFile := utils.GetConfigFile()

	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
