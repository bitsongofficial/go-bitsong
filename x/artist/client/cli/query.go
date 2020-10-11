package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the artist module",
		RunE:  client.ValidateCmd,
	}

	queryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryArtist(cdc),
	)...)

	return queryCmd
}

func GetCmdQueryArtist(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "id [artist-id]",
		Short: "query the artist by artist-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the artist by id.
Example:
$ %s query %s id 1
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]
			if id == "" {
				return nil
			}

			params := types.NewQueryArtistParams(id)
			bz := cdc.MustMarshalJSON(params)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryArtist)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("Could not find artist with id %s \n", id)
				return nil
			}

			var artist types.Artist
			cdc.MustUnmarshalJSON(res, &artist)

			return cliCtx.PrintOutput(artist)
		},
	}
}
