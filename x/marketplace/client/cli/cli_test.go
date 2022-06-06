package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	marketplacecli "github.com/bitsongofficial/go-bitsong/x/marketplace/client/cli"
	marketplacetestutil "github.com/bitsongofficial/go-bitsong/x/marketplace/client/testutil"
	marketplacetypes "github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	nfttestutil "github.com/bitsongofficial/go-bitsong/x/nft/client/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	_, err = nfttestutil.CreateNFT(clientCtx, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	_, err = marketplacetestutil.CreateAuction(clientCtx, 1, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	_, err = marketplacetestutil.StartAuction(clientCtx, 1, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	_, err = marketplacetestutil.PlaceBid(clientCtx, 1, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestGetCmdQueryAuctions() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdQueryAuctions()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp marketplacetypes.QueryAuctionsResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdQueryAuction()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp marketplacetypes.QueryAuctionResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryBidsByAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdQueryBidsByAuction()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp marketplacetypes.QueryBidsByAuctionResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryBidsByBidder() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdQueryBidsByBidder()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			val.Address.String(),
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp marketplacetypes.QueryBidsByBidderResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryBidderMetadata() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdQueryBidderMetadata()

	_, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			val.Address.String(),
			"--output=json",
		},
	)
	s.Require().Error(err)

	// var resp marketplacetypes.QueryBidderMetadataResponse
	// clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdCreateAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdCreateAuction()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagNftId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagPrizeType, "NFT_ONLY_TRANSFER"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagBidDenom, "utbsg"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagDuration, "864000s"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagPriceFloor, "1000000"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagInstantSalePrice, "100000000"),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagTickSize, "100000"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestSetAuctionAuthority() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdSetAuctionAuthority()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagNewAuthority, val.Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestStartAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdStartAuction()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestEndAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdEndAuction()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestPlaceBid() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdPlaceBid()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAmount, "100000ubtsg"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestCancelBid() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdCancelBid()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestClaimBid() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := marketplacecli.GetCmdClaimBid()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", marketplacecli.FlagAuctionId, fmt.Sprintf("%d", 1)),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}
