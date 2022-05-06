package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

// GetQueryCmd returns the query commands for the nft module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the nft module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryNFTInfo(),
		GetCmdQueryNFTsByOwner(),
		GetCmdQueryMetadata(),
		GetCmdQueryCollection(),
	)

	return queryCmd
}

func GetCmdQueryNFTInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nft-info [id]",
		Long:    "Query a nft information by id.",
		Example: fmt.Sprintf(`$ %s query nft nft-info [id]`, version.AppName),
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

			res, err := queryClient.NFTInfo(context.Background(), &types.QueryNFTInfoRequest{
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

func GetCmdQueryNFTsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nfts-by-owner [owner]",
		Long:    "Query all nfts information by owner.",
		Example: fmt.Sprintf(`$ %s query nft nfts-by-owner [owner]`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.NFTsByOwner(context.Background(), &types.QueryNFTsByOwnerRequest{
				Owner: args[0],
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

func GetCmdQueryMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metadata [id]",
		Long:    "Query a metadata information by id.",
		Example: fmt.Sprintf(`$ %s query nft metadata [id]`, version.AppName),
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

			res, err := queryClient.Metadata(context.Background(), &types.QueryMetadataRequest{
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

func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection [id]",
		Long:    "Query a collection information by id.",
		Example: fmt.Sprintf(`$ %s query nft collection [id]`, version.AppName),
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

			res, err := queryClient.Collection(context.Background(), &types.QueryCollectionRequest{
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
