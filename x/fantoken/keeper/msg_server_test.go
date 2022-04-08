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
	token := tokentypes.NewFanToken(name, sdk.NewInt(21000000), owner, denomMetaData)

	beginBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom)
	suite.Equal("100000000000000ubtsg", beginBondDenomAmt.String())

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.IssueFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgIssueFanToken(
		token.GetSymbol(), token.Name,
		token.MaxSupply, token.MetaData.Description, token.GetOwner().String(), sdk.NewCoin(types.BondDenom, sdk.NewInt(999999)),
	))
	suite.Error(err, "the issue fee is less than the standard")

	_, err = msgServer.IssueFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgIssueFanToken(
		token.GetSymbol(), token.Name,
		token.MaxSupply, token.MetaData.Description, token.GetOwner().String(), sdk.NewCoin(types.BondDenom, sdk.NewInt(1000000)),
	))
	suite.NoError(err)

	suite.True(suite.keeper.HasFanToken(suite.ctx, token.GetDenom()))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, token.GetDenom())
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwner().String())
	suite.EqualValues(&token, issuedToken.(*tokentypes.FanToken))

	endBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom)
	suite.Equal(beginBondDenomAmt.Sub(endBondDenomAmt).Amount, sdk.NewInt(1000000))
}
