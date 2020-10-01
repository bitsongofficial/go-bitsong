package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	pQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the profile module",
		RunE:  client.ValidateCmd,
	}

	pQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryCreator(cdc),
		GetCmdQueryRelease(cdc),
	)...)

	return pQueryCmd
}

func GetCmdQueryRelease(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "id [releaseID]",
		Short: "query release by releaseID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query release by releaseID.
Example:
$ %s query %s id releaseid
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			releaseID := args[0]
			if releaseID == "" {
				return nil
			}

			params := types.NewQueryReleaseParams(releaseID)
			bz := cdc.MustMarshalJSON(params)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRelease)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("Could not find release with id %s \n", releaseID)
				return nil
			}

			var release types.Release
			cdc.MustUnmarshalJSON(res, &release)

			return cliCtx.PrintOutput(release)
		},
	}
}

func GetCmdQueryCreator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "creator [creatorAddress]",
		Short: "get all releases owned by a creator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get all releases owned by a creator.
Example:
$ %s query %s creator bitsong12lmjr995d0f6dkzpplm58g5makm75eefh0n9fl
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryAllReleaseForCreatorParams(address)
			bz := cdc.MustMarshalJSON(params)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllReleaseForCreator)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("Could not find profile with address %s \n", address.String())
				return nil
			}

			var releases []types.Release
			cdc.MustUnmarshalJSON(res, &releases)

			return cliCtx.PrintOutput(releases)
		},
	}
}
