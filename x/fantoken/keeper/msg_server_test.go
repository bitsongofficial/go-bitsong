package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *KeeperTestSuite) TestMsgIssueFanToken() {
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
	token := tokentypes.NewFanToken(name, sdk.NewInt(21000000), owner, uri, denomMetaData)

	beginBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom)
	suite.Equal("100000000000000ubtsg", beginBondDenomAmt.String())

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.IssueFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgIssueFanToken(
		token.GetSymbol(), token.Name,
		token.MaxSupply, token.MetaData.Description, token.GetOwner().String(), token.GetUri(), sdk.NewCoin(types.BondDenom, sdk.NewInt(999999)),
	))
	suite.Error(err, "the issue fee is less than the standard")

	_, err = msgServer.IssueFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgIssueFanToken(
		token.GetSymbol(), token.Name,
		token.MaxSupply, token.MetaData.Description, token.GetOwner().String(), token.GetUri(), sdk.NewCoin(types.BondDenom, sdk.NewInt(1000000)),
	))
	suite.NoError(err)

	suite.True(suite.keeper.HasFanToken(suite.ctx, token.GetDenom()))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, token.GetDenom())
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwner().String())
	suite.Equal(token.URI, issuedToken.GetUri())
	suite.EqualValues(&token, issuedToken.(*tokentypes.FanToken))

	endBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom)
	suite.Equal(beginBondDenomAmt.Sub(endBondDenomAmt).Amount, sdk.NewInt(1000000))
}

func (suite *KeeperTestSuite) TestMsgEditFanToken() {
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

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.EditFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgEditFanToken(denom, mintable, owner.String()))
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

func (suite *KeeperTestSuite) TestMsgMintFanToken() {
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

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.MintFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgMintFanToken(recipient.String(), token.GetDenom(), token.GetOwner().String(), mintAmount))
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("1000ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", amt.String())

	// mint token without owner

	err = suite.keeper.MintFanToken(suite.ctx, owner, token.GetDenom(), mintAmount, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")
}

func (suite *KeeperTestSuite) TestMsgBurnFanToken() {
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

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.MintFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgMintFanToken(recipient.String(), token.GetDenom(), token.GetOwner().String(), mintAmount))
	suite.NoError(err)

	burnedAmount := sdk.NewInt(200)

	_, err = msgServer.BurnFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgBurnFanToken(token.GetDenom(), token.GetOwner().String(), burnedAmount))
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("800ft73676a7961793266743066347032627463426974636f696e204e6574776f726b", amt.String())
}
