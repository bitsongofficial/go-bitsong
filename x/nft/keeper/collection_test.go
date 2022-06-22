package keeper_test

import "github.com/bitsongofficial/go-bitsong/x/nft/types"

func (suite *KeeperTestSuite) TestLastCollectionIdGetSet() {
	// get default last collection id
	lastCollectionId := suite.app.NFTKeeper.GetLastCollectionId(suite.ctx)
	suite.Require().Equal(lastCollectionId, uint64(0))

	// set last collection id to new value
	newCollectionId := uint64(2)
	suite.app.NFTKeeper.SetLastCollectionId(suite.ctx, newCollectionId)

	// check last collection id update
	lastCollectionId = suite.app.NFTKeeper.GetLastCollectionId(suite.ctx)
	suite.Require().Equal(lastCollectionId, newCollectionId)
}

func (suite *KeeperTestSuite) TestCollectionGetSet() {
	// get collection by not available id
	_, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, 0)
	suite.Require().Error(err)

	// get all collections when not available
	allCollections := suite.app.NFTKeeper.GetAllCollections(suite.ctx)
	suite.Require().Len(allCollections, 0)

	// create a new collection
	collections := []types.Collection{
		{
			Id:              1,
			Name:            "name1",
			Uri:             "uri1",
			UpdateAuthority: "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
		},
		{
			Id:              2,
			Name:            "name2",
			Uri:             "",
			UpdateAuthority: "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
		},
		{
			Id:              3,
			Name:            "",
			Uri:             "uri2",
			UpdateAuthority: "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
		},
		{
			Id:              4,
			Name:            "",
			Uri:             "uri2",
			UpdateAuthority: "",
		},
		{
			Id:              5,
			Name:            "",
			Uri:             "",
			UpdateAuthority: "",
		},
	}

	for _, collection := range collections {
		suite.app.NFTKeeper.SetCollection(suite.ctx, collection)
	}

	for _, collection := range collections {
		c, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, collection.Id)
		suite.Require().NoError(err)
		suite.Require().Equal(collection, c)
	}

	allCollections = suite.app.NFTKeeper.GetAllCollections(suite.ctx)
	suite.Require().Len(allCollections, 5)
	suite.Require().Equal(collections, allCollections)
}

func (suite *KeeperTestSuite) TestCollectionNftsGetSet() {
	collectionId := uint64(1)

	// check nft ids by collection id query
	nftIds := suite.app.NFTKeeper.GetCollectionNftIds(suite.ctx, collectionId)
	suite.Require().Len(nftIds, 0)

	// create a new collection
	collectionRecords := []types.CollectionRecord{
		{
			NftId:        1,
			CollectionId: 1,
		},
		{
			NftId:        2,
			CollectionId: 1,
		},
		{
			NftId:        3,
			CollectionId: 1,
		},
		{
			NftId:        3,
			CollectionId: 2,
		},
		{
			NftId:        4,
			CollectionId: 2,
		},
	}

	for _, record := range collectionRecords {
		suite.app.NFTKeeper.SetCollectionNftRecord(suite.ctx, record.CollectionId, record.NftId)
	}

	// check all records
	allCollectionNftRecords := suite.app.NFTKeeper.GetAllCollectionNftRecords(suite.ctx)
	suite.Require().Len(allCollectionNftRecords, 5)
	suite.Require().Equal(collectionRecords, allCollectionNftRecords)

	// check by collection id
	collection1NftRecords := suite.app.NFTKeeper.GetCollectionNftRecords(suite.ctx, 1)
	suite.Require().Len(collection1NftRecords, 3)
	collection2NftRecords := suite.app.NFTKeeper.GetCollectionNftRecords(suite.ctx, 2)
	suite.Require().Len(collection2NftRecords, 2)

	// delete nft3 from collection2 record
	suite.app.NFTKeeper.DeleteCollectionNftRecord(suite.ctx, 2, 3)

	// check all records
	allCollectionNftRecords = suite.app.NFTKeeper.GetAllCollectionNftRecords(suite.ctx)
	suite.Require().Len(allCollectionNftRecords, 4)

	// check by collection id
	collection1NftRecords = suite.app.NFTKeeper.GetCollectionNftRecords(suite.ctx, 1)
	suite.Require().Len(collection1NftRecords, 3)
	collection2NftRecords = suite.app.NFTKeeper.GetCollectionNftRecords(suite.ctx, 2)
	suite.Require().Len(collection2NftRecords, 1)
}
