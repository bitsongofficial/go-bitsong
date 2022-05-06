package cli_test

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

func (s *IntegrationTestSuite) TestFanTokenCli() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, from)

	//------test GetCmdIssueFanToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, fantokenObj.GetSymbol()),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, fantokenObj.GetName()),
		fmt.Sprintf("--%s=%s", tokencli.FlagMaxSupply, fantokenObj.GetMaxSupply()),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, fantokenObj.GetMintable()),
		fmt.Sprintf("--%s=%s", tokencli.FlagURI, fantokenObj.GetUri()),

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
	fantokens := &[]tokentypes.FanToken{}
	bz, err = tokentestutil.QueryFanTokensExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.LegacyAmino.UnmarshalJSON(bz.Bytes(), fantokens))
	s.Require().Equal(1, len(*fantokens))

	//------test GetCmdQueryFanToken()-------------
	var fantoken *tokentypes.FanToken
	respType = proto.Message(&tokentypes.FanToken{})
	bz, err = tokentestutil.QueryFanTokenExec(clientCtx, fantokenObj.GetDenom())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))

	fantoken = respType.(*tokentypes.FanToken)
	s.Require().Equal(name, fantoken.GetName())
	s.Require().Equal(symbol, fantoken.GetSymbol())
	s.Require().Equal(uri, fantoken.GetUri())
	s.Require().Equal(from, fantoken.GetOwner())

	//------test GetCmdQueryParams()-------------
	respType = proto.Message(&tokentypes.Params{})
	bz, err = tokentestutil.QueryParamsExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	params := respType.(*tokentypes.Params)
	s.Require().NoError(err)
	expectedParams := "{\"issue_fee\":{\"denom\":\"stake\",\"amount\":\"1000000\"}}"
	result, _ := json.Marshal(params)
	s.Require().Equal(expectedParams, string(result))

	//------test GetCmdMintFanToken()-------------
	coinType := proto.Message(&sdk.Coin{})
	out, err := simapp.QueryBalanceExec(clientCtx, from.String(), fantokenObj.GetDenom())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), coinType))
	initAmount := sdk.ZeroInt()
	mintAmount := sdk.NewInt(50000000)

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", tokencli.FlagAmount, mintAmount),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	bz, err = tokentestutil.MintFanTokenExec(clientCtx, from.String(), fantokenObj.GetDenom(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	out, err = simapp.QueryBalanceExec(clientCtx, from.String(), fantokenObj.GetDenom())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), coinType))
	balance := coinType.(*sdk.Coin)
	expectedAmount := initAmount.Add(mintAmount)
	s.Require().Equal(expectedAmount, balance.Amount)

	//------test GetCmdBurnFanToken()-------------
	burnAmount := sdk.NewInt(2000000)

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagAmount, burnAmount),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	bz, err = tokentestutil.BurnFanTokenExec(clientCtx, from.String(), fantokenObj.GetDenom(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	out, err = simapp.QueryBalanceExec(clientCtx, from.String(), fantokenObj.GetDenom())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), coinType))
	balance = coinType.(*sdk.Coin)
	expectedAmount = expectedAmount.Sub(burnAmount)
	s.Require().Equal(expectedAmount, balance.Amount)

	//------test GetCmdEditFanToken()-------------
	newMintable := false

	args = []string{
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, newMintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})
	bz, err = tokentestutil.EditFanTokenExec(clientCtx, from.String(), fantokenObj.GetDenom(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	var fantoken2 *tokentypes.FanToken
	respType = proto.Message(&tokentypes.FanToken{})
	bz, err = tokentestutil.QueryFanTokenExec(clientCtx, fantoken.GetDenom())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	fantoken2 = respType.(*tokentypes.FanToken)
	s.Require().Equal(newMintable, fantoken2.GetMintable())

	//------test GetCmdTransferTokenOwner()-------------
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagRecipient, to.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	bz, err = tokentestutil.TransferFanTokenOwnerExec(clientCtx, from.String(), fantokenObj.GetDenom(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}
