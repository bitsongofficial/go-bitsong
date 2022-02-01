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
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestToken() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// ---------------------------------------------------------------------------

	from := val.Address
	symbol := "kitty"
	name := "Kitty Token"
	maxSupply := int64(200000000)
	mintable := true
	issueFee := "1000000ubtsg"
	description := "Kitty Token"
	baseURL := val.APIAddress
	denom := tokentypes.GetFantokenDenom(from, symbol, name)

	//------test GetCmdIssueFanToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),
		fmt.Sprintf("--%s=%s", tokencli.FlagIssueFee, issueFee),
		fmt.Sprintf("--%s=%s", tokencli.FlagDescription, description),

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
	tokensResp := respType.(*tokentypes.QueryFanTokensResponse)
	s.Require().Equal(1, len(tokensResp.Tokens))

	//------test GetCmdQueryFanToken()-------------
	url = fmt.Sprintf("%s/bitsong/fantoken/v1beta1/denom/%s", baseURL, denom)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryFanTokenResponse{})
	var token tokentypes.FanTokenI
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	tokenResp := respType.(*tokentypes.QueryFanTokenResponse)
	token = tokenResp.Token
	s.Require().NoError(err)
	s.Require().Equal(name, token.GetName())
	s.Require().Equal(symbol, token.GetSymbol())

	//------test GetCmdQueryParams()-------------
	url = fmt.Sprintf("%s/bitsong/fantoken/v1beta1/params", baseURL)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryParamsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	paramsResp := respType.(*tokentypes.QueryParamsResponse)
	s.Require().NoError(err)
	expectedParams := "{\"issue_price\":{\"denom\":\"ubtsg\",\"amount\":\"1000000\"}}"
	result, _ := json.Marshal(paramsResp.Params)
	s.Require().Equal(expectedParams, string(result))
}
