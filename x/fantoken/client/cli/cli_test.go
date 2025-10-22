package cli_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	tokencli "github.com/bitsongofficial/go-bitsong/x/fantoken/client/cli"
)

var (
	name      = "Bitcoin"
	symbol    = "btc"
	uri       = "ipfs://"
	maxSupply = math.NewInt(200000000)
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.DefaultConfig()
	cfg.NumValidators = 2

	s.cfg = cfg
	s.network, _ = network.New(s.T(), "test_data", cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

// / NOTE(hard-nett):We can have confidence this modules cli functionality with the ict test package,
// which are just dockerized wrapper of the image making cli calls themselves.
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func issueCmd(s *IntegrationTestSuite, ctx client.Context, name, symbol string, maxSupply math.Int, uri string, from sdk.AccAddress) {
	//------test GetCmdIssue()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%s", tokencli.FlagURI, uri),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdIssue(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

}

func (s *IntegrationTestSuite) TestCmdIssue() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)
}
