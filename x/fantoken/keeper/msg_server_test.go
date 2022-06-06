package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestMsgIssueFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)

	beginBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, sdk.DefaultBondDenom)
	suite.Equal("100000000000000stake", beginBondDenomAmt.String())

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.IssueFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgIssueFanToken(
		fantokenObj.GetName(), fantokenObj.GetSymbol(), fantokenObj.GetUri(),
		fantokenObj.GetMaxSupply(), fantokenObj.GetOwner().String(),
	))
	suite.NoError(err)

	fts := suite.keeper.GetFanTokens(suite.ctx, owner)
	ft := fts[0]

	suite.True(suite.keeper.HasFanToken(suite.ctx, ft.GetDenom()))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, ft.GetDenom())
	suite.NoError(err)

	suite.Equal(fantokenObj.Owner, issuedToken.GetOwner().String())
	suite.Equal(fantokenObj.MetaData.URI, issuedToken.GetUri())

	endBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, sdk.DefaultBondDenom)
	suite.Equal(beginBondDenomAmt.Sub(endBondDenomAmt).Amount, sdk.NewInt(1000000))
}

func (suite *KeeperTestSuite) TestMsgEditFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.setFanToken(fantokenObj)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.EditFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgEditFanToken(fantokenObj.GetDenom(), false, owner.String()))
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, fantokenObj.GetDenom())
	suite.NoError(err)

	fantokenObj.Mintable = false
	fantokenObj.MaxSupply = sdk.ZeroInt()
	suite.EqualValues(newToken.(*tokentypes.FanToken), &fantokenObj)
}

func (suite *KeeperTestSuite) TestMsgMintFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	amt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("0ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.MintFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgMintFanToken(recipient.String(), fantokenObj.GetDenom(), fantokenObj.GetOwner().String(), mintAmount))
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("1000ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())

	// mint token without owner

	err = suite.keeper.MintFanToken(suite.ctx, owner, fantokenObj.GetDenom(), mintAmount, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")
}

func (suite *KeeperTestSuite) TestMsgBurnFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	amt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("0ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	msgServer := keeper.NewMsgServerImpl(suite.keeper)
	_, err := msgServer.MintFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgMintFanToken(recipient.String(), fantokenObj.GetDenom(), fantokenObj.GetOwner().String(), mintAmount))
	suite.NoError(err)

	burnedAmount := sdk.NewInt(200)

	_, err = msgServer.BurnFanToken(sdk.WrapSDKContext(suite.ctx), tokentypes.NewMsgBurnFanToken(fantokenObj.GetDenom(), fantokenObj.GetOwner().String(), burnedAmount))
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("800ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())
}
