package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (suite *KeeperTestSuite) TestQueryCollection() {
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

	res, err := suite.keeper.Collection(suite.ctx, &types.QueryCollectionRequest{
		Collection: collectionDenom,
	})
	suite.NoError(err)
	suite.Equal(testCollection1.Name, res.Collection.Name)
	suite.Equal(testCollection1.Symbol, res.Collection.Symbol)
	suite.Equal(testCollection1.Uri, res.Collection.Uri)
	suite.Equal(testCollection1.Minter, res.Collection.Minter)
}

func (suite *KeeperTestSuite) TestQueryOwnerOf() {
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

	res, err := suite.keeper.OwnerOf(suite.ctx, &types.QueryOwnerOfRequest{
		Collection: collectionDenom,
		TokenId:    "1",
	})
	suite.NoError(err)
	suite.Equal(owner1.String(), res.Owner)
}

func (suite *KeeperTestSuite) TestQueryNftInfo() {
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

	res, err := suite.keeper.NftInfo(suite.ctx, &types.QueryNftInfoRequest{
		Collection: collectionDenom,
		TokenId:    testNft1.TokenId,
	})
	suite.NoError(err)
	suite.Equal(testNft1.TokenId, res.Nft.TokenId)
	suite.Equal(testNft1.Name, res.Nft.Name)
	suite.Equal(testNft1.Uri, res.Nft.Uri)
	suite.Equal(collectionDenom, res.Nft.Collection)
	suite.Equal(owner1.String(), res.Nft.Owner)
}

func (suite *KeeperTestSuite) TestQueryNftsOfOwner() {
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

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom,
		testNft2.TokenId,
		testNft2.Name,
		testNft2.Uri,
	)
	suite.NoError(err)

	res, err := suite.keeper.Nfts(suite.ctx, &types.QueryNftsRequest{
		Collection: collectionDenom,
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 2)
	suite.Equal(testNft1.TokenId, res.Nfts[0].TokenId)
	suite.Equal(testNft2.TokenId, res.Nfts[1].TokenId)
}

func (suite *KeeperTestSuite) TestQueryNftsByOwner() {
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

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom,
		testNft2.TokenId,
		testNft2.Name,
		testNft2.Uri,
	)
	suite.NoError(err)

	res, err := suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner1.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 2)
	suite.Equal(testNft1.TokenId, res.Nfts[0].TokenId)
	suite.Equal(uint64(0), res.Nfts[0].Editions)
	suite.Equal(testNft2.TokenId, res.Nfts[1].TokenId)
	suite.Equal(uint64(0), res.Nfts[1].Editions)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 0)

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter1,
		owner2,
		collectionDenom,
		testNft3.TokenId,
		testNft3.Name,
		testNft3.Uri,
	)
	suite.NoError(err)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 1)
	suite.Equal(testNft3.TokenId, res.Nfts[0].TokenId)
}
