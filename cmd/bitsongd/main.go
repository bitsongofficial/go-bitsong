package main

import (
	"os"

	"github.com/bitsongofficial/bitsong/app"
	"github.com/bitsongofficial/bitsong/cmd/bitsongd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
