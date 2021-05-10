package rest_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	simapp "github.com/bitsongofficial/bitsong/app"
	tokencli "github.com/bitsongofficial/bitsong/x/fantoken/client/cli"
	tokentestutil "github.com/bitsongofficial/bitsong/x/fantoken/client/testutil"
	tokentypes "github.com/bitsongofficial/bitsong/x/fantoken/types"
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
	denom := "kitty"
	name := "Kitty Token"
	maxSupply := int64(200000000)
	mintable := true
	baseURL := val.APIAddress

	//------test GetCmdIssueFanToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagDenom, denom),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	bz, err := tokentestutil.IssueFanTokenExec(clientCtx, from.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
	tokenDenom := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()

	//------test GetCmdQueryFanTokens()-------------
	url := fmt.Sprintf("%s/bitsong/fantoken/tokens", baseURL)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryFanTokensResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(resp, respType))
	tokensResp := respType.(*tokentypes.QueryFanTokensResponse)
	s.Require().Equal(2, len(tokensResp.FanTokens))

	//------test GetCmdQueryFanToken()-------------
	url = fmt.Sprintf("%s/bitsong/fantoken/tokens/%s", baseURL, tokenDenom)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryFanTokenResponse{})
	var token tokentypes.FanTokenI
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(resp, respType))
	tokenResp := respType.(*tokentypes.QueryFanTokenResponse)
	err = clientCtx.InterfaceRegistry.UnpackAny(tokenResp.FanToken, &token)
	s.Require().NoError(err)
	s.Require().Equal(name, token.GetName())
	s.Require().Equal(denom, token.GetDenom())

	//------test GetCmdQueryParams()-------------
	url = fmt.Sprintf("%s/bitsong/fantoken/params", baseURL)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&tokentypes.QueryParamsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(resp, respType))
	paramsResp := respType.(*tokentypes.QueryParamsResponse)
	s.Require().NoError(err)
	expectedParams := "{\"issue_price\":{\"denom\":\"stake\",\"amount\":\"60000\"}}"
	result, _ := json.Marshal(paramsResp.Params)
	s.Require().Equal(expectedParams, string(result))
}
