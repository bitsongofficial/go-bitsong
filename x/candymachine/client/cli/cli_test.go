package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	candymachinecli "github.com/bitsongofficial/go-bitsong/x/candymachine/client/cli"
	candymachinetypes "github.com/bitsongofficial/go-bitsong/x/candymachine/types"
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

func (s *IntegrationTestSuite) GetCmdQueryCandyMachines() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := candymachinecli.GetCmdQueryCandyMachines()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp candymachinetypes.QueryCandyMachinesResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryCandyMachine() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := candymachinecli.GetCmdQueryCandyMachine()

	_, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	// var resp candymachinetypes.QueryCandyMachineResponse
	// clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryBidsByAuction() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := candymachinecli.GetCmdQueryParams()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp candymachinetypes.QueryParamsResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdCreateCandyMachine() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := candymachinecli.GetCmdCreateCandyMachine()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagCollectionId, 1),
		fmt.Sprintf("--%s=1000", candymachinecli.FlagPrice),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagDenom, "utbsg"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagMetadataBaseUrl, "https://punk.com/metadata"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagEndSettingsType, "BY_MINT"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagEndSettingsValue, "10"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagTreasury, val.Address.String()),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagGoLiveDate, "1659404536"),
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

func (s *IntegrationTestSuite) TestGetCmdUpdateCandyMachine() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := candymachinecli.GetCmdUpdateCandyMachine()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagCollectionId, 1),
		fmt.Sprintf("--%s=1000", candymachinecli.FlagPrice),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagDenom, "utbsg"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagMetadataBaseUrl, "https://punk.com/metadata"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagEndSettingsType, "BY_MINT"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagEndSettingsValue, "10"),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagTreasury, val.Address.String()),
		fmt.Sprintf("--%s=%s", candymachinecli.FlagGoLiveDate, "1659404536"),
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

	cmd := candymachinecli.GetCmdMintNFT()

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

func (s *IntegrationTestSuite) TestGetCmdCloseCandyMachine() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := candymachinecli.GetCmdCloseCandyMachine()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("1"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}
