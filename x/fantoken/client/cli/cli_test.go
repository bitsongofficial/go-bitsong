package cli_test

import (
	"fmt"
	simapp "github.com/bitsongofficial/go-bitsong/app"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/cli"
	"testing"

	tokencli "github.com/bitsongofficial/go-bitsong/x/fantoken/client/cli"
)

var (
	name      = "Bitcoin"
	symbol    = "btc"
	uri       = "ipfs://"
	maxSupply = sdk.NewInt(200000000)
	mintable  = true
	height    = int64(1)
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
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestFanTokenCli() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	//------test GetCmdIssue()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%s", tokencli.FlagURI, uri),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdIssue(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	denom := string(txResp.Events[12].Attributes[0].Value)
	denom = denom[1 : len(denom)-1]

	//------test GetCmdQueryFanToken()-------------
	var result tokentypes.QueryFanTokenResponse

	args = []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &result))
	s.Require().Equal(result.Fantoken.Denom, denom)

	//------test GetCmdQueryFanTokens()-------------
	var results tokentypes.QueryFanTokensResponse

	args = []string{
		from.String(),
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err = clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanTokens(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &results))
	s.Require().Len(results.Fantokens, 1)
}
