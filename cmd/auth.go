package cmd

import (
	"fmt"
	"os"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/kwngo/hop-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(configureCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Set up auth",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure api credentials for strapi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Configure auth credentials...")
		validate := func(input string) error {
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "Set your strapi api key",
			Validate: validate,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			fmt.Println(err)
			os.Exit(1)
		}
		hostPrompt := promptui.Prompt{
			Label:    "Set your strapi host url",
			Validate: validate,
		}

		hostResult, err := hostPrompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			fmt.Println(err)
			os.Exit(1)
		}

		cfg := HopConfig{
			AuthToken: result,
			Host:      hostResult,
		}
		bytes, err := toml.Marshal(cfg)
		if err != nil {
			fmt.Println("There was an issue unmarshaling auth credentials.")
			os.Exit(1)
		}
		cfgFile := utils.GetConfigFile()
		f, err := os.OpenFile(cfgFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("There was an issue configuring auth credentials.")
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()
		f.WriteString(string(bytes))
	},
}
