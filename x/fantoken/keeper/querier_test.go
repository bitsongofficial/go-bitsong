package keeper_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/bitsongofficial/ledger/x/fantoken/keeper"
	"github.com/bitsongofficial/ledger/x/fantoken/types"
)

func (suite *KeeperTestSuite) TestQueryFanToken() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.legacyAmino)

	params := types.QueryFanTokenParams{
		Denom: types.GetNativeToken().Denom,
	}
	bz := suite.legacyAmino.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryFanToken),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryFanToken}, query)
	suite.Nil(err)

	data2 := codec.MustMarshalJSONIndent(suite.legacyAmino, types.GetNativeToken())
	suite.Equal(data2, data)
}

func (suite *KeeperTestSuite) TestQueryFanTokens() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.legacyAmino)

	params := types.QueryFanTokensParams{
		Owner: nil,
	}
	bz := suite.legacyAmino.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryFanTokens),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryFanTokens}, query)
	suite.Nil(err)

	data2 := codec.MustMarshalJSONIndent(suite.legacyAmino, []types.FanTokenI{types.GetNativeToken()})
	suite.Equal(data2, data)
}
