package main

import (
	"os"

	"cosmossdk.io/log"
	"github.com/bitsongofficial/go-bitsong/v018/app/params"

	"github.com/bitsongofficial/go-bitsong/v018/app"
	"github.com/bitsongofficial/go-bitsong/v018/cmd/bitsongd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	params.SetAddressPrefixes()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "BITSONGD", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
