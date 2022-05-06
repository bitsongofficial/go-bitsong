package keeper_test

import "github.com/bitsongofficial/go-bitsong/x/nft/types"

func (suite *KeeperTestSuite) TestLastMetadataIdGetSet() {
	// get default last metadata id
	lastMetadataId := suite.app.NFTKeeper.GetLastMetadataId(suite.ctx)
	suite.Require().Equal(lastMetadataId, uint64(0))

	// set last metadata id to new value
	newMetadataId := uint64(2)
	suite.app.NFTKeeper.SetLastMetadataId(suite.ctx, newMetadataId)

	// check last metadata id update
	lastMetadataId = suite.app.NFTKeeper.GetLastMetadataId(suite.ctx)
	suite.Require().Equal(lastMetadataId, newMetadataId)
}

func (suite *KeeperTestSuite) TestMetadataGetSet() {
	// get metadata by not available id
	_, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, 0)
	suite.Require().Error(err)

	// get all metadata when not available
	allMetadata := suite.app.NFTKeeper.GetAllMetadata(suite.ctx)
	suite.Require().Len(allMetadata, 0)

	// create new metadata
	data := &types.Data{
		Name:                 "meta1",
		Symbol:               "META1",
		Uri:                  "uri1",
		SellerFeeBasisPoints: 10,
		Creators: []*types.Creator{
			{
				Address:  "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
				Verified: false,
				Share:    1,
			},
		},
	}
	metadata := []types.Metadata{
		{
			Id:                  1,
			UpdateAuthority:     "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Mint:                "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Data:                data,
			PrimarySaleHappened: false,
			IsMutable:           true,
		},
		{
			Id:                  2,
			UpdateAuthority:     "",
			Mint:                "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Data:                data,
			PrimarySaleHappened: false,
			IsMutable:           true,
		},
		{
			Id:                  3,
			UpdateAuthority:     "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Mint:                "",
			Data:                data,
			PrimarySaleHappened: false,
			IsMutable:           true,
		},
		{
			Id:                  4,
			UpdateAuthority:     "",
			Mint:                "",
			Data:                data,
			PrimarySaleHappened: false,
			IsMutable:           true,
		},
	}

	for _, meta := range metadata {
		suite.app.NFTKeeper.SetMetadata(suite.ctx, meta)
	}

	for _, meta := range metadata {
		m, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, meta.Id)
		suite.Require().NoError(err)
		suite.Require().Equal(meta, m)
	}

	allMetadata = suite.app.NFTKeeper.GetAllMetadata(suite.ctx)
	suite.Require().Len(allMetadata, 4)
	suite.Require().Equal(metadata, allMetadata)
}
