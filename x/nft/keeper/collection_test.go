package keeper_test

func (suite *KeeperTestSuite) TestCollectionIdGetSet() {

	lastCollectionId := suite.app.NFTKeeper.GetLastCollectionId(suite.ctx)
	suite.Require().Equal(lastCollectionId, uint64(0))

	newCollectionId := uint64(2)
	suite.app.NFTKeeper.SetLastCollectionId(suite.ctx, newCollectionId)

	lastCollectionId = suite.app.NFTKeeper.GetLastCollectionId(suite.ctx)
	suite.Require().Equal(lastCollectionId, newCollectionId)
}

// TODO: test
// GetCollectionById
// SetCollection
// GetAllCollections

// TODO: test
// SetCollectionNftRecord
// DeleteCollectionNftRecord
// GetCollectionNftRecords
// GetAllCollectionNftRecords
