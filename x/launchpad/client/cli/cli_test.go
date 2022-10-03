package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	launchpadcli "github.com/bitsongofficial/go-bitsong/x/launchpad/client/cli"
	launchpadtypes "github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	nftcli "github.com/bitsongofficial/go-bitsong/x/nft/client/cli"
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

	cfg := simapp.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	_, err = nfttestutil.CreateCollection(clientCtx, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	_, err = nfttestutil.CreateNFT(clientCtx, val.Address.String(), s.cfg.BondDenom)
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

func (s *IntegrationTestSuite) GetCmdQueryLaunchPads() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdQueryLaunchPads()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp launchpadtypes.QueryLaunchPadsResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryLaunchPad() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdQueryLaunchPad()

	_, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	// var resp launchpadtypes.QueryLaunchPadResponse
	// clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryBidsByAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdQueryParams()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp launchpadtypes.QueryParamsResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdCreateLaunchPad() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdCreateLaunchPad()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagCollectionId, 1),
		fmt.Sprintf("--%s=1000", launchpadcli.FlagPrice),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagDenom, "utbsg"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagMetadataBaseUrl, "https://punk.com/metadata"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagEndTimestamp, "0"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagMaxMint, "10"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagTreasury, val.Address.String()),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagGoLiveDate, "1659404536"),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreators, val.Address.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreatorShares, "10"),
		fmt.Sprintf("--%s=true", nftcli.FlagMutable),
		fmt.Sprintf("--%s=%d", nftcli.FlagSellerFeeBasisPoints, 100),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdUpdateLaunchPad() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdUpdateLaunchPad()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagCollectionId, 1),
		fmt.Sprintf("--%s=1000", launchpadcli.FlagPrice),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagDenom, "utbsg"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagMetadataBaseUrl, "https://punk.com/metadata"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagEndTimestamp, "0"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagMaxMint, "10"),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagTreasury, val.Address.String()),
		fmt.Sprintf("--%s=%s", launchpadcli.FlagGoLiveDate, "1659404536"),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreators, val.Address.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreatorShares, "10"),
		fmt.Sprintf("--%s=true", nftcli.FlagMutable),
		fmt.Sprintf("--%s=%d", nftcli.FlagSellerFeeBasisPoints, 100),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdMintNFT() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdMintNFT()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("1"),
		fmt.Sprintf("MyPunk"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdCloseLaunchPad() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := launchpadcli.GetCmdCloseLaunchPad()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("1"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}
