package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	nftcli "github.com/bitsongofficial/go-bitsong/x/nft/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/nft/client/testutil"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
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
	_, err = testutil.CreateCollection(clientCtx, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = testutil.CreateNFT(clientCtx, val.Address.String(), s.cfg.BondDenom)
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

func (s *IntegrationTestSuite) TestGetCmdQueryNFTInfo() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdQueryNFTInfo()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1:1:0",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp nfttypes.QueryNFTInfoResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryMetadata() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdQueryMetadata()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			"1",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp nfttypes.QueryMetadataResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdQueryCollection() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdQueryCollection()

	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			"--output=json",
		},
	)
	s.Require().NoError(err)

	var resp nfttypes.QueryCollectionResponse
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &resp)
}

func (s *IntegrationTestSuite) TestGetCmdCreateNFT() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdCreateNFT()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagName, "Punk10"),
		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, "PUNK"),
		fmt.Sprintf("--%s=%s", nftcli.FlagUri, "https://punk.com/10"),
		fmt.Sprintf("--%s=%d", nftcli.FlagSellerFeeBasisPoints, 100),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreators, val.Address.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreatorShares, "10"),
		fmt.Sprintf("--%s=false", nftcli.FlagMutable),
		fmt.Sprintf("--%s=%s", nftcli.FlagUpdateAuthority, val.Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdTransferNFT() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdTransferNFT()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagNftId, "1:1:0"),
		fmt.Sprintf("--%s=%s", nftcli.FlagNewOwner, val.Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdSignMetadata() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdSignMetadata()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagMetadataId, 1),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdUpdateMetadata() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdUpdateMetadata()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagMetadataId, 1),
		fmt.Sprintf("--%s=%s", nftcli.FlagName, "Punk10"),
		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, "PUNK"),
		fmt.Sprintf("--%s=%s", nftcli.FlagUri, "https://punk.com/10"),
		fmt.Sprintf("--%s=%d", nftcli.FlagSellerFeeBasisPoints, 100),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreators, val.Address.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagCreatorShares, "10"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdUpdateMetadataAuthority() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdUpdateMetadataAuthority()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagMetadataId, 1),
		fmt.Sprintf("--%s=%s", nftcli.FlagNewAuthority, val.Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdCreateCollection() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdCreateCollection()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagName, "Punk"),
		fmt.Sprintf("--%s=%s", nftcli.FlagUri, "https://punk.com"),
		fmt.Sprintf("--%s=%s", nftcli.FlagUpdateAuthority, val.Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetCmdUpdateCollectionAuthority() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := nftcli.GetCmdUpdateCollectionAuthority()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", nftcli.FlagCollectionId, 1),
		fmt.Sprintf("--%s=%s", nftcli.FlagNewAuthority, val.Address.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}
