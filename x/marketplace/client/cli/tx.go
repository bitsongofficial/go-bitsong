package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// NewTxCmd returns the transaction commands for the marketplace module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "marketplace transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateAuction(),
		GetCmdSetAuctionAuthority(),
		GetCmdStartAuction(),
		GetCmdEndAuction(),
		GetCmdPlaceBid(),
		GetCmdCancelBid(),
		GetCmdClaimBid(),
	)

	return txCmd
}

func GetCmdCreateAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-auction [flags]",
		Long: "Create auction from provided params",
		Example: fmt.Sprintf(
			`$ %s tx marketplace create-auction
				--nft-id=1
				--prize-type="NFT_ONLY_TRANSFER"
				--bid-denom="ubtsg"
				--duration="864000s"
				--price-floor=1000000
				--instant-sale-price=100000000
				--tick-size=100000`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			nftId, err := cmd.Flags().GetUint64(FlagNftId)
			if err != nil {
				return err
			}
			prizeTypeStr, err := cmd.Flags().GetString(FlagPrizeType)
			if err != nil {
				return err
			}
			bidDenom, err := cmd.Flags().GetString(FlagBidDenom)
			if err != nil {
				return err
			}
			duration, err := cmd.Flags().GetDuration(FlagDuration)
			if err != nil {
				return err
			}
			priceFloor, err := cmd.Flags().GetUint64(FlagPriceFloor)
			if err != nil {
				return err
			}
			instantSalePrice, err := cmd.Flags().GetUint64(FlagInstantSalePrice)
			if err != nil {
				return err
			}
			tickSize, err := cmd.Flags().GetUint64(FlagTickSize)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAuction(
				clientCtx.GetFromAddress(), nftId, types.AuctionPrizeType(types.AuctionPrizeType_value[prizeTypeStr]),
				bidDenom, duration, priceFloor, instantSalePrice, tickSize,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateAuction())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdSetAuctionAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set-auction-authority [flags]",
		Long: "Set auction authority to a new address",
		Example: fmt.Sprintf(
			`$ %s tx marketplace set-auction-authority
				--auction-id=1
				--new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47"`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := cmd.Flags().GetUint64(FlagAuctionId)
			if err != nil {
				return err
			}
			newAuthority, err := cmd.Flags().GetString(FlagNewAuthority)
			if err != nil {
				return err
			}
			msg := types.NewMsgSetAuctionAuthority(
				clientCtx.GetFromAddress(), auctionId, newAuthority,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetAuctionAuthority())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdStartAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "start-auction [flags]",
		Long: "Start auction by authority",
		Example: fmt.Sprintf(
			`$ %s tx marketplace start-auction
				--auction-id=1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := cmd.Flags().GetUint64(FlagAuctionId)
			if err != nil {
				return err
			}
			msg := types.NewMsgStartAuction(
				clientCtx.GetFromAddress(), auctionId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagStartAuction())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdEndAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "end-auction [flags]",
		Long: "End auction by authority",
		Example: fmt.Sprintf(
			`$ %s tx marketplace end-auction
				--auction-id=1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := cmd.Flags().GetUint64(FlagAuctionId)
			if err != nil {
				return err
			}
			msg := types.NewMsgEndAuction(
				clientCtx.GetFromAddress(), auctionId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagEndAuction())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdPlaceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "place-bid [flags]",
		Long: "Place bid on an auction",
		Example: fmt.Sprintf(
			`$ %s tx marketplace place-bid
				--auction-id=1
				--amount="1000000ubtsg"`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := cmd.Flags().GetUint64(FlagAuctionId)
			if err != nil {
				return err
			}

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(amountStr)
			if err != nil {
				return fmt.Errorf("failed to parse bid amount: %s", amountStr)
			}

			msg := types.NewMsgPlaceBid(
				clientCtx.GetFromAddress(), auctionId, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagPlaceBid())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdCancelBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "cancel-bid [flags]",
		Long: "cancel bid on an auction",
		Example: fmt.Sprintf(
			`$ %s tx marketplace cancel-bid
				--auction-id=1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := cmd.Flags().GetUint64(FlagAuctionId)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelBid(
				clientCtx.GetFromAddress(), auctionId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCancelBid())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdClaimBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "claim-bid [flags]",
		Long: "claim bid on an auction",
		Example: fmt.Sprintf(
			`$ %s tx marketplace claim-bid
				--auction-id=1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := cmd.Flags().GetUint64(FlagAuctionId)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimBid(
				clientCtx.GetFromAddress(), auctionId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagClaimBid())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
