package utils

import (
	"errors"
	"fmt"
	"os"
)

var cfgFile string

func GetConfigFile() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	path := home + "/.hopcli.toml"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()
	}
	cfgFile = path
	return cfgFile
}
