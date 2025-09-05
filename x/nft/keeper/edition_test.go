package keeper_test

func (suite *KeeperTestSuite) TestPrintEdition() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		"",
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom,
		testNft1.TokenId,
		testNft1.Name,
		testNft1.Uri,
	)
	suite.NoError(err)

	edition, err := suite.keeper.PrintEdition(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom,
		"1",
	)
	suite.NoError(err)
	suite.Equal(uint64(1), edition)

	edition, err = suite.keeper.PrintEdition(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom,
		"1",
	)
	suite.NoError(err)
	suite.Equal(uint64(2), edition)

	edition, err = suite.keeper.PrintEdition(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom,
		"1",
	)
	suite.NoError(err)
	suite.Equal(uint64(3), edition)
}
