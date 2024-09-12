package main

import (
	"os"

	"github.com/bitsongofficial/go-bitsong/app/params"

	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/cmd/bitsongd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	params.SetAddressPrefixes()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "BITSONGD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
