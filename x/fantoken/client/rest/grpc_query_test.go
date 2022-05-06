package rest_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	tokencli "github.com/bitsongofficial/go-bitsong/x/fantoken/client/cli"
	tokentestutil "github.com/bitsongofficial/go-bitsong/x/fantoken/client/testutil"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

var (
	name      = "Bitcoin"
	symbol    = "btc"
	uri       = "ipfs://"
	maxSupply = sdk.NewInt(200000000)
	mintable  = true
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

func (s *IntegrationTestSuite) TestFanToken() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address
	baseURL := val.APIAddress
	// ---------------------------------------------------------------------------

	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, from)

	//------test GetCmdIssueFanToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, fantokenObj.GetSymbol()),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, fantokenObj.GetName()),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, fantokenObj.GetMaxSupply().Int64()),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, fantokenObj.GetMintable()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	bz, err := tokentestutil.IssueFanTokenExec(clientCtx, from.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQueryFanTokens()-------------
	url := fmt.Sprintf("%s/bitsong/fantoken/v1beta1/fantokens", baseURL)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryFanTokensResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	fantokensResp := respType.(*tokentypes.QueryFanTokensResponse)
	s.Require().Equal(1, len(fantokensResp.Fantokens))

	//------test GetCmdQueryFanToken()-------------
	url = fmt.Sprintf("%s/bitsong/fantoken/v1beta1/denom/%s", baseURL, fantokenObj.GetDenom())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryFanTokenResponse{})
	var fantoken tokentypes.FanTokenI
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	fantokenResp := respType.(*tokentypes.QueryFanTokenResponse)
	fantoken = fantokenResp.Fantoken
	s.Require().NoError(err)
	s.Require().Equal(fantokenObj.GetName(), fantoken.GetName())
	s.Require().Equal(fantokenObj.GetSymbol(), fantoken.GetSymbol())

	//------test GetCmdQueryParams()-------------
	url = fmt.Sprintf("%s/bitsong/fantoken/v1beta1/params", baseURL)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryParamsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	paramsResp := respType.(*tokentypes.QueryParamsResponse)
	s.Require().NoError(err)
	expectedParams := "{\"issue_fee\":{\"denom\":\"stake\",\"amount\":\"1000000\"}}"
	result, _ := json.Marshal(paramsResp.Params)
	s.Require().Equal(expectedParams, string(result))
}
