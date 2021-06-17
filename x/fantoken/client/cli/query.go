package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
)

// GetQueryCmd returns the query commands for the token module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the token module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryFanToken(),
		GetCmdQueryFanTokens(),
		GetCmdQueryTotalBurn(),
		GetCmdQueryParams(),
	)

	return queryCmd
}

// GetCmdQueryToken implements the query token command.
func GetCmdQueryFanToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom]",
		Long:    "Query a token by symbol or denom.",
		Example: fmt.Sprintf("$ %s query token token <denom>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			if err := types.ValidateSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FanToken(context.Background(), &types.QueryFanTokenRequest{
				Denom: args[0],
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.FanToken)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokens implements the query tokens command.
func GetCmdQueryFanTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tokens [owner]",
		Long:    "Query token by the owner.",
		Example: fmt.Sprintf("$ %s query token tokens <owner>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.FanTokens(
				context.Background(),
				&types.QueryFanTokensRequest{
					Owner:      owner.String(),
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			tokens := make([]types.FanTokenI, 0, len(res.FanTokens))
			for _, eviAny := range res.FanTokens {
				var evi types.FanTokenI
				if err = clientCtx.InterfaceRegistry.UnpackAny(eviAny, &evi); err != nil {
					return err
				}
				tokens = append(tokens, evi)
			}

			return clientCtx.PrintObjectLegacy(tokens)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all tokens")

	return cmd
}

// GetCmdQueryParams implements the query token related param command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Long:    "Query values set as token parameters.",
		Example: fmt.Sprintf("$ %s query token params", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTotalBurn return the total amount of all burned tokens
func GetCmdQueryTotalBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "total-burn",
		Long:    "Query the total amount of all burned tokens.",
		Example: fmt.Sprintf("$ %s query token params", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TotalBurn(context.Background(), &types.QueryTotalBurnRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
