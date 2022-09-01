package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestParamsGetSet() {
	params := suite.app.LaunchPadKeeper.GetParamSet(suite.ctx)
	params.LaunchpadCreationPrice = sdk.NewInt64Coin("ubtsg", 1000000)
	params.LaunchpadMaxMint = 1000

	suite.app.LaunchPadKeeper.SetParamSet(suite.ctx, params)
	newParams := suite.app.LaunchPadKeeper.GetParamSet(suite.ctx)
	suite.Require().Equal(params, newParams)
}
