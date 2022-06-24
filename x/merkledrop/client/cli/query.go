package cli

import (
	"context"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strconv"
)

// GetQueryCmd returns the query commands for the nft module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the merkledrop module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryMerkledrop(),
		GetCmdQueryIndexClaimed(),
	)

	return queryCmd
}

func GetCmdQueryMerkledrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "detail [id]",
		Long:    "Query a merkledrop detail by id.",
		Example: fmt.Sprintf(`$ %s query merkledrop detail [id]`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Merkledrop(context.Background(), &types.QueryMerkledropRequest{
				Id: uint64(id),
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

func GetCmdQueryIndexClaimed() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index-claimed [id] [index]",
		Long:    "Query if an index and id have been claimed.",
		Example: fmt.Sprintf(`$ %s query merkledrop index-claimed [id] [index]`, version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			index, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IndexClaimed(context.Background(), &types.QueryIndexClaimedRequest{
				Id:    uint64(id),
				Index: uint64(index),
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
