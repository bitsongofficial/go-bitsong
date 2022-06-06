package fantoken_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	tokenmodule "github.com/bitsongofficial/go-bitsong/x/fantoken"
	tokenkeeper "github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

const (
	isCheckTx = false
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite

	cdc    codec.JSONCodec
	ctx    sdk.Context
	keeper tokenkeeper.Keeper
	bk     bankkeeper.Keeper
}

func (suite *HandlerSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.FanTokenKeeper
	suite.bk = app.BankKeeper

	// set params
	suite.keeper.SetParamSet(suite.ctx, tokentypes.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, tokentypes.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, tokentypes.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func (suite *HandlerSuite) issueFanToken(fantoken tokentypes.FanToken) {
	err := suite.keeper.AddFanToken(suite.ctx, fantoken)
	suite.NoError(err)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			fantoken.GetDenom(),
			sdk.ZeroInt(),
		),
	)

	err = suite.bk.MintCoins(suite.ctx, tokentypes.ModuleName, mintCoins)
	suite.NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, tokentypes.ModuleName, owner, mintCoins)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestIssueFanToken() {
	h := tokenmodule.NewHandler(suite.keeper)

	nativeTokenAmt1 := suite.bk.GetBalance(suite.ctx, owner, sdk.DefaultBondDenom).Amount

	msg := tokentypes.NewMsgIssueFanToken(name, symbol, uri, maxSupply, owner.String())

	res, err := h(suite.ctx, msg)
	suite.NoError(err)

	denom := string(res.Events[4].Attributes[0].Value)
	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)
	suite.Equal(uri, issuedToken.GetUri())

	nativeTokenAmt2 := suite.bk.GetBalance(suite.ctx, owner, sdk.DefaultBondDenom).Amount

	suite.Equal(nativeTokenAmt1.Sub(issueFee.Amount), nativeTokenAmt2)

	nativeTokenAmt3 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	suite.Equal(nativeTokenAmt3, sdk.ZeroInt())
}

func (suite *HandlerSuite) TestMintFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom()).Amount
	suite.Equal(sdk.ZeroInt(), beginBtcAmt)

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintFanToken := tokentypes.NewMsgMintFanToken("", fantokenObj.GetDenom(), fantokenObj.Owner, sdk.NewInt(1000))
	_, err := h(suite.ctx, msgMintFanToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom()).Amount
	mintBtcAmt := msgMintFanToken.Amount
	suite.Equal(beginBtcAmt.Add(mintBtcAmt), endBtcAmt)
}

func (suite *HandlerSuite) TestBurnFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintFanToken := tokentypes.NewMsgMintFanToken("", fantokenObj.GetDenom(), fantokenObj.Owner, sdk.NewInt(1000))
	_, err := h(suite.ctx, msgMintFanToken)
	suite.NoError(err)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom()).Amount
	suite.Equal(sdk.NewInt(1000), beginBtcAmt)

	msgBurnFanToken := tokentypes.NewMsgBurnFanToken(fantokenObj.GetDenom(), fantokenObj.Owner, sdk.NewInt(200))
	_, err = h(suite.ctx, msgBurnFanToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom()).Amount
	burnBtcAmt := msgBurnFanToken.Amount

	suite.Equal(beginBtcAmt.Sub(burnBtcAmt), endBtcAmt)
}
