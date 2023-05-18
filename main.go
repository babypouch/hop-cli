package main

import (
	"github.com/babypouch/hop-cli/cmd"
	_ "github.com/babypouch/hop-cli/cmd/mentions"
	_ "github.com/babypouch/hop-cli/cmd/products"
)

func main() {
	cmd.Execute()
}
