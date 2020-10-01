package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/profile/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		GetCmdQueryAddress(cdc),
		GetCmdQueryProfile(cdc),
	)...)

	return pQueryCmd
}

func GetCmdQueryProfile(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "handle [handle]",
		Short: "query the profile by handle",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the profile by handle.
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

			params := types.NewQueryProfileParams(handle)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProfile)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return sdkerrors.Wrap(types.ErrProfileNotFound, fmt.Sprintf("%s", handle))
			}

			var profile types.Profile
			err = cdc.UnmarshalJSON(res, &profile)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(profile)
		},
	}
}

func GetCmdQueryAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address [accountAddress]",
		Short: "get the profile owned by an account address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the profile owned by an account address.
Example:
$ %s query %s owner bitsong12lmjr995d0f6dkzpplm58g5makm75eefh0n9fl
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

			params := types.NewQueryByAddressParams(address)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProfileByAddress)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return sdkerrors.Wrap(types.ErrProfileNotFound, fmt.Sprintf("%s", address.String()))
			}

			var profile types.Profile
			err = cdc.UnmarshalJSON(res, &profile)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(profile)
		},
	}
}
