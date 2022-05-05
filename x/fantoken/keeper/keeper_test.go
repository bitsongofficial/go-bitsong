package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	isCheckTx = false
)

var (
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	uri      = "ipfs://"
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initAmt)}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	bk          bankkeeper.Keeper
	app         *simapp.BitsongApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.FanTokenKeeper
	suite.bk = app.BankKeeper
	suite.app = app

	// set params
	suite.keeper.SetParamSet(suite.ctx, tokentypes.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, tokentypes.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, tokentypes.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setFanToken(token tokentypes.FanToken) {
	err := suite.keeper.AddFanToken(suite.ctx, token)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) issueFanToken(token tokentypes.FanToken) {
	suite.setFanToken(token)
}

func (suite *KeeperTestSuite) TestIssueFanToken() {
	symbol := "btc"
	name := "Bitcoin Network"
	denom := tokentypes.GetFantokenDenom(owner, symbol, name)
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        denom,
		Display:     symbol,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: denom, Exponent: 0},
			{Denom: symbol, Exponent: tokentypes.FanTokenDecimal},
		},
	}
	uri := "ipfs://"
	token := tokentypes.NewFanToken(name, sdk.NewInt(21000000), owner, uri, denomMetaData)

	_, err := suite.keeper.IssueFanToken(
		suite.ctx, token.GetSymbol(), token.Name,
		token.MaxSupply, token.MetaData.Description, token.GetOwner(), token.GetUri(), sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(999999)),
	)
	suite.Error(err, "the issue fee is less than the standard")

	_, err = suite.keeper.IssueFanToken(
		suite.ctx, token.GetSymbol(), token.Name,
		token.MaxSupply, token.MetaData.Description, token.GetOwner(), token.GetUri(), sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)),
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasFanToken(suite.ctx, token.GetDenom()))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, token.GetDenom())
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwner().String())
	suite.Equal(token.URI, issuedToken.GetUri())

	suite.EqualValues(&token, issuedToken.(*tokentypes.FanToken))
}

func (suite *KeeperTestSuite) TestEditFanToken() {
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(21000000), owner, uri, denomMetaData)
	suite.setFanToken(token)

	denom := "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b"
	mintable := false

	err := suite.keeper.EditFanToken(suite.ctx, denom, mintable, owner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)

	expToken := tokentypes.FanToken{
		Name:      "Bitcoin Network",
		MaxSupply: sdk.ZeroInt(),
		Mintable:  false,
		Owner:     owner.String(),
		URI:       uri,
		MetaData:  denomMetaData,
	}

	suite.EqualValues(newToken.(*tokentypes.FanToken), &expToken)
}

func (suite *KeeperTestSuite) TestMintFanToken() {
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(2000), owner, uri, denomMetaData)
	suite.issueFanToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("0ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, token.GetDenom(), mintAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("1000ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", amt.String())

	// mint token without owner

	err = suite.keeper.MintFanToken(suite.ctx, owner, token.GetDenom(), mintAmount, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")
}

func (suite *KeeperTestSuite) TestBurnFanToken() {
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(2000), owner, uri, denomMetaData)
	suite.issueFanToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("0ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, token.GetDenom(), mintAmount, token.GetOwner())
	suite.NoError(err)

	burnedAmount := sdk.NewInt(200)

	err = suite.keeper.BurnFanToken(suite.ctx, token.GetDenom(), burnedAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("800ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", amt.String())
}

func (suite *KeeperTestSuite) TestTransferFanToken() {
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(21000000), owner, uri, denomMetaData)
	suite.setFanToken(token)

	dstOwner := sdk.AccAddress(tmhash.SumTruncated([]byte("TokenDstOwner")))

	err := suite.keeper.TransferFanTokenOwner(suite.ctx, token.GetDenom(), token.GetOwner(), dstOwner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, token.GetDenom())
	suite.NoError(err)

	suite.Equal(dstOwner, newToken.GetOwner())
}
