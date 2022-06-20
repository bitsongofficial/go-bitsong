package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// GetQueryCmd returns the query commands for the fantoken module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the fantoken module",
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

// GetCmdQueryFanToken implements the query fantoken command.
func GetCmdQueryFanToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom [denom]",
		Long:    "Query a fantoken by denom.",
		Example: fmt.Sprintf("$ %s query fantoken denom <denom>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			if err := types.ValidateDenom(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FanToken(context.Background(), &types.QueryFanTokenRequest{
				Denom: args[0],
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryFanTokens implements the query fantokens command.
func GetCmdQueryFanTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [owner]",
		Long:    "Query fantokens by the owner.",
		Example: fmt.Sprintf("$ %s query fantoken owner <owner>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			authority, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.FanTokens(context.Background(), &types.QueryFanTokensRequest{
				Authority:  authority.String(),
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all fantokens")

	return cmd
}

// GetCmdQueryParams implements the query fantoken related param command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Long:    "Query values set as fantoken parameters.",
		Example: fmt.Sprintf("$ %s query fantoken params", version.AppName),
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

// GetCmdQueryTotalBurn return the total amount of all burned fantokens
func GetCmdQueryTotalBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "total-burn",
		Long:    "Query the total amount of all burned fantokens.",
		Example: fmt.Sprintf("$ %s query fantoken total-burn", version.AppName),
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
