package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/channel/types"
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
	queryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the channel module",
		RunE:  client.ValidateCmd,
	}

	queryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryByOwner(cdc),
		GetCmdQueryChannel(cdc),
	)...)

	return queryCmd
}

func GetCmdQueryChannel(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "handle [handle]",
		Short: "query the channel by handle",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the channel by handle.
Example:
$ %s query %s handle test
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			handle := args[0]
			if handle == "" {
				return nil
			}

			params := types.NewQueryChannelParams(handle)
			bz := cdc.MustMarshalJSON(params)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChannel)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("Could not find channel with handle %s \n", handle)
				return nil
			}

			var channel types.Channel
			cdc.MustUnmarshalJSON(res, &channel)

			return cliCtx.PrintOutput(channel)
		},
	}
}

func GetCmdQueryByOwner(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "owner [accountAddress]",
		Short: "get the channel owned by an account address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the channel owned by an account address.
Example:
$ %s query %s owner bitsong12lmjr995d0f6dkzpplm58g5makm75eefh0n9fl
`, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryByOwnerParams(owner)
			bz := cdc.MustMarshalJSON(params)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChannelByOwner)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("Could not find channel with address %s \n", owner.String())
				return nil
			}

			var channel types.Channel
			cdc.MustUnmarshalJSON(res, &channel)

			return cliCtx.PrintOutput(channel)
		},
	}
}
