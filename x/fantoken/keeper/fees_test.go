package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/chainmodules/types"
	tokentypes "github.com/bitsongofficial/chainmodules/x/fantoken/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *KeeperTestSuite) TestDeductIssueFanTokenFee() {
	beginBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom)
	suite.Equal("100000000000000ubtsg", beginBondDenomAmt.String())

	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ubtc",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ubtc", Exponent: 0},
			{Denom: "btc", Exponent: tokentypes.FanTokenDecimal},
		},
	}
	token := tokentypes.NewFanToken("Bitcoin Network", sdk.NewInt(2000), owner, denomMetaData)
	suite.issueFanToken(token)

	issueFeeAmt := sdk.NewInt(1000000)
	err := suite.keeper.DeductIssueFanTokenFee(suite.ctx, owner, sdk.NewCoin(types.BondDenom, issueFeeAmt), token.GetSymbol())
	suite.NoError(err)

	endBondDenomAmt := suite.bk.GetBalance(suite.ctx, owner, types.BondDenom)
	suite.Equal(beginBondDenomAmt.Sub(endBondDenomAmt).Amount, issueFeeAmt)
}
