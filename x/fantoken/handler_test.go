package token_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	simapp "github.com/bitsongofficial/ledger/app"
	tokenmodule "github.com/bitsongofficial/ledger/x/fantoken"
	tokenkeeper "github.com/bitsongofficial/ledger/x/fantoken/keeper"
	"github.com/bitsongofficial/ledger/x/fantoken/types"
)

const (
	isCheckTx = false
)

var (
	nativeToken = types.GetNativeToken()
	denom       = nativeToken.Denom
	owner       = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmt     = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin    = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper tokenkeeper.Keeper
	bk     bankkeeper.Keeper
}

func (suite *HandlerSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func (suite *HandlerSuite) issueFanToken(token types.FanToken) {
	err := suite.keeper.AddFanToken(suite.ctx, token)
	suite.NoError(err)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.Denom,
			sdk.ZeroInt(),
		),
	)

	err = suite.bk.MintCoins(suite.ctx, types.ModuleName, mintCoins)
	suite.NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, mintCoins)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestIssueFanToken() {
	h := tokenmodule.NewHandler(suite.keeper)

	nativeTokenAmt1 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	msg := types.NewMsgIssueFanToken("btc", "satoshi", sdk.NewInt(21000000), false, owner.String())

	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	nativeTokenAmt2 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	params := suite.keeper.GetParamSet(suite.ctx)

	suite.Equal(nativeTokenAmt1.Sub(params.IssuePrice.Amount), nativeTokenAmt2)

	nativeTokenAmt3 := suite.bk.GetBalance(suite.ctx, owner, msg.Denom).Amount
	suite.Equal(nativeTokenAmt3, sdk.ZeroInt())
}

func (suite *HandlerSuite) TestMintFanToken() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(2000), true, owner)
	suite.issueFanToken(token)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom).Amount
	suite.Equal(sdk.ZeroInt(), beginBtcAmt)

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintFanToken := types.NewMsgMintFanToken("", token.Denom, token.Owner, sdk.NewInt(1000))
	_, err := h(suite.ctx, msgMintFanToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom).Amount
	mintBtcAmt := msgMintFanToken.Amount
	suite.Equal(beginBtcAmt.Add(mintBtcAmt), endBtcAmt)
}

func (suite *HandlerSuite) TestBurnFanToken() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(2000), true, owner)
	suite.issueFanToken(token)

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintFanToken := types.NewMsgMintFanToken("", token.Denom, token.Owner, sdk.NewInt(1000))
	_, err := h(suite.ctx, msgMintFanToken)
	suite.NoError(err)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom).Amount
	suite.Equal(sdk.NewInt(1000), beginBtcAmt)

	msgBurnFanToken := types.NewMsgBurnFanToken(token.Denom, token.Owner, sdk.NewInt(200))
	_, err = h(suite.ctx, msgBurnFanToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom).Amount
	burnBtcAmt := msgBurnFanToken.Amount

	suite.Equal(beginBtcAmt.Sub(burnBtcAmt), endBtcAmt)
}
