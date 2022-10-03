package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
)

// GetQueryCmd returns the query commands for the launchpad module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the launchpad module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryLaunchPads(),
		GetCmdQueryLaunchPad(),
		GetCmdQueryMintableMetadataIds(),
	)

	return queryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "params [flags]",
		Long: "Query params.",
		Example: fmt.Sprintf(
			`$ %s query launchpad params`, version.AppName),
		Args: cobra.ExactArgs(0),
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

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryLaunchPads() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "launchpads [flags]",
		Long: "Query launchpads.",
		Example: fmt.Sprintf(
			`$ %s query launchpad launchpads`, version.AppName),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LaunchPads(context.Background(), &types.QueryLaunchPadsRequest{})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryLaunchPad() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "launchpad [collection_id] [flags]",
		Long: "Query a launchpad by collection id.",
		Example: fmt.Sprintf(
			`$ %s query launchpad launchpad 1`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			collId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.LaunchPad(context.Background(), &types.QueryLaunchPadRequest{
				CollId: uint64(collId),
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

func GetCmdQueryMintableMetadataIds() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mintabe-metadata-ids [collection_id] [flags]",
		Long: "Query mintable metadata ids by collection id.",
		Example: fmt.Sprintf(
			`$ %s query launchpad mintabe-metadata-ids 1`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			collId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.MintableMetadataIds(context.Background(), &types.QueryMintableMetadataIdsRequest{
				CollId: uint64(collId),
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Info)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
