package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (suite *KeeperTestSuite) TestQueryCollection() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, collectionDenom)

	res, err := suite.keeper.Collection(suite.ctx, &types.QueryCollectionRequest{
		Collection: collectionDenom,
	})
	suite.NoError(err)
	suite.Equal(testCollection.Name, res.Collection.Name)
	suite.Equal(testCollection.Symbol, res.Collection.Symbol)
	suite.Equal(testCollection.Description, res.Collection.Description)
	suite.Equal(testCollection.Uri, res.Collection.Uri)
	suite.Equal(testCollection.Minter, res.Collection.Minter)
}

func (suite *KeeperTestSuite) TestQueryOwnerOf() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, collectionDenom)

	nft1 := types.Nft{
		TokenId:     "1",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft1)
	suite.NoError(err)

	res, err := suite.keeper.OwnerOf(suite.ctx, &types.QueryOwnerOfRequest{
		Collection: collectionDenom,
		TokenId:    "1",
	})
	suite.NoError(err)
	suite.Equal(owner.String(), res.Owner)
}

func (suite *KeeperTestSuite) TestQueryNumTokens() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, collectionDenom)

	supply := suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(uint64(0), supply.Uint64())

	nft1 := types.Nft{
		TokenId:     "1",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	nft2 := types.Nft{
		TokenId:     "2",
		Name:        "My Second NFT",
		Description: "This is my second NFT",
		Uri:         "ipfs://my-second-nft-metadata.json",
	}

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft1)
	suite.NoError(err)

	supply = suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(uint64(1), supply.Uint64())

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft2)
	suite.NoError(err)

	supply = suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(uint64(2), supply.Uint64())

	res, err := suite.keeper.NumTokens(suite.ctx, &types.QueryNumTokensRequest{
		Collection: collectionDenom,
	})
	suite.NoError(err)
	suite.Equal(uint64(2), res.Count)
}

func (suite *KeeperTestSuite) TestQueryNftInfo() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, collectionDenom)

	nft1 := types.Nft{
		TokenId:     "1",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft1)
	suite.NoError(err)

	res, err := suite.keeper.NftInfo(suite.ctx, &types.QueryNftInfoRequest{
		Collection: collectionDenom,
		TokenId:    "1",
	})
	suite.NoError(err)
	suite.Equal(nft1.TokenId, res.Nft.TokenId)
	suite.Equal(nft1.Name, res.Nft.Name)
	suite.Equal(nft1.Description, res.Nft.Description)
	suite.Equal(nft1.Uri, res.Nft.Uri)
	suite.Equal(collectionDenom, res.Nft.Collection)
	suite.Equal(owner.String(), res.Nft.Owner)
}

func (suite *KeeperTestSuite) TestQueryNftsOfOwner() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, collectionDenom)

	nft1 := types.Nft{
		TokenId:     "1",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	nft2 := types.Nft{
		TokenId:     "2",
		Name:        "My Second NFT",
		Description: "This is my second NFT",
		Uri:         "ipfs://my-second-nft-metadata.json",
	}

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft1)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft2)
	suite.NoError(err)

	res, err := suite.keeper.Nfts(suite.ctx, &types.QueryNftsRequest{
		Collection: collectionDenom,
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 2)
	suite.Equal(nft1.TokenId, res.Nfts[0].TokenId)
	suite.Equal(nft2.TokenId, res.Nfts[1].TokenId)
}

func (suite *KeeperTestSuite) TestQueryNftsByOwner() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, collectionDenom)

	nft1 := types.Nft{
		TokenId:     "1",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	nft2 := types.Nft{
		TokenId:     "2",
		Name:        "My Second NFT",
		Description: "This is my second NFT",
		Uri:         "ipfs://my-second-nft-metadata.json",
	}

	nft3 := types.Nft{
		TokenId:     "3",
		Name:        "My Third NFT",
		Description: "This is my third NFT",
		Uri:         "ipfs://my-third-nft-metadata.json",
	}

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft1)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft2)
	suite.NoError(err)

	res, err := suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 2)
	suite.Equal(nft1.TokenId, res.Nfts[0].TokenId)
	suite.Equal(nft2.TokenId, res.Nfts[1].TokenId)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 0)

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner2, nft3)
	suite.NoError(err)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 1)
	suite.Equal(nft3.TokenId, res.Nfts[0].TokenId)
}
