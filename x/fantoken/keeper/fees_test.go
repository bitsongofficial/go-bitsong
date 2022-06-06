package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

func (suite *KeeperTestSuite) TestDeductIssueFanTokenFee() {
	beginBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, sdk.DefaultBondDenom)
	suite.Equal("100000000000000stake", beginBondDenomAmt.String())

	height := int64(1)
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	issueFeeAmt := sdk.NewInt(1000000)
	err := suite.keeper.DeductIssueFanTokenFee(suite.ctx, owner)
	suite.NoError(err)

	endBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, sdk.DefaultBondDenom)
	suite.Equal(beginBondDenomAmt.Sub(endBondDenomAmt).Amount, issueFeeAmt)
}
