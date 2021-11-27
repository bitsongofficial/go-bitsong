package cli_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/bitsongofficial/chainmodules/app"
	tokencli "github.com/bitsongofficial/chainmodules/x/fantoken/client/cli"
	tokentestutil "github.com/bitsongofficial/chainmodules/x/fantoken/client/testutil"
	tokentypes "github.com/bitsongofficial/chainmodules/x/fantoken/types"
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
	denom := tokentypes.GetFantokenDenom(from, symbol, name)
	maxSupply := sdk.NewInt(200000000)
	mintable := true
	issueFee := "1000000ubtsg"
	description := "Kitty Token"
	//------test GetCmdIssueFanToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMaxSupply, maxSupply),
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
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQueryFanTokens()-------------
	tokens := &[]tokentypes.FanToken{}
	bz, err = tokentestutil.QueryFanTokensExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.LegacyAmino.UnmarshalJSON(bz.Bytes(), tokens))
	s.Require().Equal(1, len(*tokens))

	//------test GetCmdQueryFanToken()-------------
	var token *tokentypes.FanToken
	respType = proto.Message(&tokentypes.FanToken{})
	bz, err = tokentestutil.QueryFanTokenExec(clientCtx, denom)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	token = respType.(*tokentypes.FanToken)
	s.Require().Equal(name, token.GetName())
	s.Require().Equal(symbol, token.GetSymbol())

	//------test GetCmdQueryParams()-------------
	respType = proto.Message(&tokentypes.Params{})
	bz, err = tokentestutil.QueryParamsExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	params := respType.(*tokentypes.Params)
	s.Require().NoError(err)
	expectedParams := "{\"issue_price\":{\"denom\":\"ubtsg\",\"amount\":\"1000000\"}}"
	result, _ := json.Marshal(params)
	s.Require().Equal(expectedParams, string(result))

	//------test GetCmdMintFanToken()-------------
	coinType := proto.Message(&sdk.Coin{})
	out, err := simapp.QueryBalanceExec(clientCtx, from.String(), symbol)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), coinType))
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
	bz, err = tokentestutil.MintFanTokenExec(clientCtx, from.String(), denom, args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	out, err = simapp.QueryBalanceExec(clientCtx, from.String(), denom)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), coinType))
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
	bz, err = tokentestutil.BurnFanTokenExec(clientCtx, from.String(), denom, args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	out, err = simapp.QueryBalanceExec(clientCtx, from.String(), denom)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), coinType))
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
	bz, err = tokentestutil.EditFanTokenExec(clientCtx, from.String(), denom, args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	var token2 *tokentypes.FanToken
	respType = proto.Message(&tokentypes.FanToken{})
	bz, err = tokentestutil.QueryFanTokenExec(clientCtx, denom)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	token2 = respType.(*tokentypes.FanToken)
	s.Require().Equal(newMintable, token2.GetMintable())

	//------test GetCmdTransferTokenOwner()-------------
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagRecipient, to.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	bz, err = tokentestutil.TransferFanTokenOwnerExec(clientCtx, from.String(), denom, args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}
