package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	marketplacecli "github.com/bitsongofficial/go-bitsong/x/marketplace/client/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

func CreateAuction(clientCtx client.Context, nftId uint64, from string, bondDenom string) (testutil.BufferWriter, error) {
	cmd := marketplacecli.GetCmdCreateAuction()

	return clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagNftId, fmt.Sprintf("%d", nftId)),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagPrizeType, "NFT_ONLY_TRANSFER"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagBidDenom, "utbsg"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagDuration, "864000s"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagPriceFloor, "1000000"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagInstantSalePrice, "100000000"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagTickSize, "100000"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(100))).String()),
	})
}

func StartAuction(clientCtx client.Context, auctionId uint64, from string, bondDenom string) (testutil.BufferWriter, error) {
	cmd := marketplacecli.GetCmdStartAuction()

	return clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", auctionId)),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(100))).String()),
	})
}

func PlaceBid(clientCtx client.Context, auctionId uint64, from string, bondDenom string) (testutil.BufferWriter, error) {
	cmd := marketplacecli.GetCmdPlaceBid()

	return clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", auctionId)),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAmount, "10000000ubtsg"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(100))).String()),
	})
}
