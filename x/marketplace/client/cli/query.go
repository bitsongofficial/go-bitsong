package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
)

// GetQueryCmd returns the query commands for the marketplace module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the marketplace module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryAuctions(),
		GetCmdQueryAuction(),
		GetCmdQueryBidsByAuction(),
		GetCmdQueryBidsByBidder(),
		GetCmdQueryBidderMetadata(),
	)

	return queryCmd
}

func GetCmdQueryAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "auctions [flags]",
		Long: "Query auctions by auction state and authority.",
		Example: fmt.Sprintf(
			`$ %s query marketplace auctions
				--state="EMPTY"
				--authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47"`, version.AppName),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			authority, err := cmd.Flags().GetString(FlagAuthority)
			if err != nil {
				return err
			}

			stateStr, err := cmd.Flags().GetString(FlagAuctionState)
			if err != nil {
				return err
			}

			state := types.AuctionState_value[stateStr]

			res, err := queryClient.Auctions(context.Background(), &types.QueryAuctionsRequest{
				Authority: authority,
				State:     types.AuctionState(state),
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(FlagQueryAuctions())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auction [id]",
		Long:    "Query auction by id.",
		Example: fmt.Sprintf(`$ %s query marketplace auction 1`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Auction(context.Background(), &types.QueryAuctionRequest{
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

func GetCmdQueryBidsByAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bids-by-auction [id]",
		Long:    "Query bids by auction id.",
		Example: fmt.Sprintf(`$ %s query marketplace bids-by-auction 1`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.BidsByAuction(context.Background(), &types.QueryBidsByAuctionRequest{
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

func GetCmdQueryBidsByBidder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bids-by-bidder [bidder]",
		Long:    "Query bids by bidder.",
		Example: fmt.Sprintf(`$ %s query marketplace bids-by-bidder bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BidsByBidder(context.Background(), &types.QueryBidsByBidderRequest{
				Bidder: args[0],
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

func GetCmdQueryBidderMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bidder-metadata [bidder]",
		Long:    "Query bidder metadata by bidder.",
		Example: fmt.Sprintf(`$ %s query marketplace bidder-metadata bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BidderMetadata(context.Background(), &types.QueryBidderMetadataRequest{
				Bidder: args[0],
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
