package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestParamsGetSet() {
	params := suite.app.CandyMachineKeeper.GetParamSet(suite.ctx)
	params.CandymachineCreationPrice = sdk.NewInt64Coin("ubtsg", 1000000)
	params.CandymachineMaxMint = 1000

	suite.app.CandyMachineKeeper.SetParamSet(suite.ctx, params)
	newParams := suite.app.CandyMachineKeeper.GetParamSet(suite.ctx)
	suite.Require().Equal(params, newParams)
}
