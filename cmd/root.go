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

// var cfgFile string

// func init() {
// 	cobra.OnInitialize(initConfig)
// 	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hopcli)")

// }

// func initConfig() {
// 	// Don't forget to read config either from cfgFile or from home directory!
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		// Search config in home directory with name ".cobra" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".cobra")
// 	}

// 	if err := viper.ReadInConfig(); err != nil {
// 		fmt.Println("Can't read config:", err)
// 		os.Exit(1)
// 	}
// }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
