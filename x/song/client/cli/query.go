package cli

import (
	"fmt"

	"github.com/BitSongOfficial/go-bitsong/x/song/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QuerySongsParams struct {
	Owner sdk.AccAddress
}

func NewQuerySongsParams(addr sdk.AccAddress) QuerySongsParams {
	return QuerySongsParams{
		Owner: addr,
	}
}

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	songQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the song module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	songQueryCmd.AddCommand(client.GetCommands(
		GetCmdList(cdc),
	)...)
	return songQueryCmd
}

func GetCmdList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "list [address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Search all songs of a specific address",
		Example: "$ bitsongcli query song list bitsong1hf4n743fujvxrwx8af7u35anjqpdd2cx8p6cdd",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := NewQuerySongsParams(owner)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/songs", types.QuerierRoute), bz)
			if err != nil {
				return err
			}

			var songs types.Songs
			cdc.MustUnmarshalJSON(res, &songs)
			return cliCtx.PrintOutput(songs)
		},
	}
}
