package cli

import (
	"github.com/spf13/cobra"

	"github.com/BitSongOfficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	songTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Artist transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	songTxCmd.AddCommand(client.PostCommands(
		GetCmdRegister(cdc),
	)...)

	return songTxCmd
}

// GetCmdRegister is the CLI command for register an Artist transaction
func GetCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "register",
		Short:   "Register a new artist",
		Example: "$ bitsongcli artist register --name=\"Armin Van Buuren\" --image=<ipfs_url> --from mykey",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
