package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
)

// GetQueryCmd returns the query commands for the candymachine module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the candymachine module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryCandyMachines(),
		GetCmdQueryCandyMachine(),
		GetCmdQueryMintableMetadataIds(),
	)

	return queryCmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "params [flags]",
		Long: "Query params.",
		Example: fmt.Sprintf(
			`$ %s query candymachine params`, version.AppName),
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

func GetCmdQueryCandyMachines() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "candymachines [flags]",
		Long: "Query candymachines.",
		Example: fmt.Sprintf(
			`$ %s query candymachine candymachines`, version.AppName),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CandyMachines(context.Background(), &types.QueryCandyMachinesRequest{})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryCandyMachine() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "candymachine [collection_id] [flags]",
		Long: "Query a candymachine by collection id.",
		Example: fmt.Sprintf(
			`$ %s query candymachine candymachine 1`, version.AppName),
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

			res, err := queryClient.CandyMachine(context.Background(), &types.QueryCandyMachineRequest{
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
			`$ %s query candymachine mintabe-metadata-ids 1`, version.AppName),
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
