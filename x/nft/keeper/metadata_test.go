package keeper_test

import "github.com/bitsongofficial/go-bitsong/x/nft/types"

func (suite *KeeperTestSuite) TestLastMetadataIdGetSet() {
	// get default last metadata id
	lastMetadataId := suite.app.NFTKeeper.GetLastMetadataId(suite.ctx, 1)
	suite.Require().Equal(lastMetadataId, uint64(0))

	// set last metadata id to new value
	newMetadataId := uint64(2)
	suite.app.NFTKeeper.SetLastMetadataId(suite.ctx, 1, newMetadataId)

	// check last metadata id update
	lastMetadataId = suite.app.NFTKeeper.GetLastMetadataId(suite.ctx, 1)
	suite.Require().Equal(lastMetadataId, newMetadataId)
}

func (suite *KeeperTestSuite) TestMetadataGetSet() {
	// get metadata by not available id
	_, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, 0, 0)
	suite.Require().Error(err)

	// get all metadata when not available
	allMetadata := suite.app.NFTKeeper.GetAllMetadata(suite.ctx)
	suite.Require().Len(allMetadata, 0)

	// create new metadata
	creators := []types.Creator{
		{
			Address:  "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Verified: false,
			Share:    1,
		},
	}
	metadata := []types.Metadata{
		{
			CollId:               1,
			Id:                   1,
			MetadataAuthority:    "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			MintAuthority:        "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Name:                 "meta1",
			Uri:                  "uri1",
			SellerFeeBasisPoints: 10,
			Creators:             creators,
			PrimarySaleHappened:  false,
			IsMutable:            true,
		},
		{
			CollId:               1,
			Id:                   2,
			MetadataAuthority:    "",
			MintAuthority:        "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			Name:                 "meta1",
			Uri:                  "uri1",
			SellerFeeBasisPoints: 10,
			Creators:             creators,
			PrimarySaleHappened:  false,
			IsMutable:            true,
		},
		{
			CollId:               1,
			Id:                   3,
			MetadataAuthority:    "bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47",
			MintAuthority:        "",
			Name:                 "meta1",
			Uri:                  "uri1",
			SellerFeeBasisPoints: 10,
			Creators:             creators,
			PrimarySaleHappened:  false,
			IsMutable:            true,
		},
		{
			CollId:               1,
			Id:                   4,
			MetadataAuthority:    "",
			MintAuthority:        "",
			Name:                 "meta1",
			Uri:                  "uri1",
			SellerFeeBasisPoints: 10,
			Creators:             creators,
			PrimarySaleHappened:  false,
			IsMutable:            true,
		},
	}

	for _, meta := range metadata {
		suite.app.NFTKeeper.SetMetadata(suite.ctx, meta)
	}

	for _, meta := range metadata {
		m, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, meta.CollId, meta.Id)
		suite.Require().NoError(err)
		suite.Require().Equal(meta, m)
	}

	allMetadata = suite.app.NFTKeeper.GetAllMetadata(suite.ctx)
	suite.Require().Len(allMetadata, 4)
	suite.Require().Equal(metadata, allMetadata)
}
