package fantoken_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	tokenmodule "github.com/bitsongofficial/go-bitsong/x/fantoken"
	tokenkeeper "github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bitsongofficial/go-bitsong/types"
)

const (
	isCheckTx = false
)

var (
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	uri      = "ipfs://"
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(types.BondDenom, initAmt)}
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

func (suite *HandlerSuite) issueFanToken(token tokentypes.FanToken) {
	err := suite.keeper.AddFanToken(suite.ctx, token)
	suite.NoError(err)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.GetDenom(),
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

	symbol := "btc"
	name := "satoshi"
	issueFee := sdk.NewCoin(types.BondDenom, sdk.NewInt(1000000))

	nativeTokenAmt1 := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom).Amount

	msg := tokentypes.NewMsgIssueFanToken(symbol, name, sdk.NewInt(21000000), "test", owner.String(), uri, issueFee)
	denom := tokentypes.GetFantokenDenom(owner, msg.Symbol, msg.Name)

	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)
	suite.Equal(uri, issuedToken.GetUri())

	nativeTokenAmt2 := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom).Amount

	suite.Equal(nativeTokenAmt1.Sub(issueFee.Amount), nativeTokenAmt2)

	nativeTokenAmt3 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	suite.Equal(nativeTokenAmt3, sdk.ZeroInt())
}

func (suite *HandlerSuite) TestMintFanToken() {
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ftbtc",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ftbtc", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(2000), owner, uri, denomMetaData)
	suite.issueFanToken(token)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom()).Amount
	suite.Equal(sdk.ZeroInt(), beginBtcAmt)

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintFanToken := tokentypes.NewMsgMintFanToken("", token.GetDenom(), token.Owner, sdk.NewInt(1000))
	_, err := h(suite.ctx, msgMintFanToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom()).Amount
	mintBtcAmt := msgMintFanToken.Amount
	suite.Equal(beginBtcAmt.Add(mintBtcAmt), endBtcAmt)
}

func (suite *HandlerSuite) TestBurnFanToken() {
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ftbtc",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ftbtc", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(2000), owner, uri, denomMetaData)
	suite.issueFanToken(token)

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintFanToken := tokentypes.NewMsgMintFanToken("", token.GetDenom(), token.Owner, sdk.NewInt(1000))
	_, err := h(suite.ctx, msgMintFanToken)
	suite.NoError(err)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom()).Amount
	suite.Equal(sdk.NewInt(1000), beginBtcAmt)

	msgBurnFanToken := tokentypes.NewMsgBurnFanToken(token.GetDenom(), token.Owner, sdk.NewInt(200))
	_, err = h(suite.ctx, msgBurnFanToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom()).Amount
	burnBtcAmt := msgBurnFanToken.Amount

	suite.Equal(beginBtcAmt.Sub(burnBtcAmt), endBtcAmt)
}
