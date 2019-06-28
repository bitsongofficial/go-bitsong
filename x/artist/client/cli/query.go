package cli

import (
	"github.com/BitSongOfficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	artistQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the artist module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	artistQueryCmd.AddCommand(client.GetCommands(
		GetCmdList(cdc),
	)...)
	return artistQueryCmd
}

func GetCmdList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "list [address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Search all artists of a specific address",
		Example: "$ bitsongcli query artist list bitsong1hf4n743fujvxrwx8af7u35anjqpdd2cx8p6cdd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
