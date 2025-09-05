package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	simapp "github.com/bitsongofficial/go-bitsong/app"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

var (
	creator1 = sdk.AccAddress(tmhash.SumTruncated([]byte("creator1")))
	creator2 = sdk.AccAddress(tmhash.SumTruncated([]byte("creator2")))

	minter1 = sdk.AccAddress(tmhash.SumTruncated([]byte("minter1")))
	minter2 = sdk.AccAddress(tmhash.SumTruncated([]byte("minter2")))

	owner1 = sdk.AccAddress(tmhash.SumTruncated([]byte("owner1")))
	owner2 = sdk.AccAddress(tmhash.SumTruncated([]byte("owner2")))

	testCollection1 = types.Collection{
		Name:   "My NFT Collection",
		Symbol: "MYNFT",
		Uri:    "ipfs://my-nft-collection-metadata.json",
		Minter: minter1.String(),
	}
	expectedDenom1 = "nft9436DDD23FB751AEA7BC6C767F20F943DD735E06"

	testNft1 = types.Nft{
		TokenId: "1",
		Name:    "My First NFT",
		Uri:     "ipfs://my-first-nft-metadata.json",
	}

	testNft2 = types.Nft{
		TokenId: "2",
		Name:    "My Second NFT",
		Uri:     "ipfs://my-second-nft-metadata.json",
	}

	testNft3 = types.Nft{
		TokenId: "3",
		Name:    "My Third NFT",
		Uri:     "ipfs://my-third-nft-metadata.json",
	}
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	ctx    sdk.Context
	keeper keeper.Keeper
	app    *simapp.BitsongApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	app := suite.App
	suite.keeper = app.NftKeeper
	suite.App = app
	suite.ctx = suite.Ctx
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestCreateCollection() {
	denom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1,
		minter1,
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom1, denom)

	_, err = suite.keeper.CreateCollection(
		suite.ctx,
		creator1,
		minter1,
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.Error(err)
}

func (suite *KeeperTestSuite) TestMintNFT() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1,
		minter1,
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom1, collectionDenom)

	supply := suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(math.NewInt(0), supply)

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

	supply = suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(math.NewInt(1), supply)

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

	supply = suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(math.NewInt(2), supply)
}

func (suite *KeeperTestSuite) TestSendNFT() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1,
		minter1,
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

	res, err := suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner1.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 1)
	suite.Equal(testNft1.TokenId, res.Nfts[0].TokenId)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 0)

	err = suite.keeper.SendNft(suite.ctx, owner1, owner2, collectionDenom, "1")
	suite.NoError(err)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 1)
	suite.Equal(testNft1.TokenId, res.Nfts[0].TokenId)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner1.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 0)

	nft, err := suite.keeper.NFTs.Get(suite.ctx, collections.Join(collectionDenom, "1"))
	suite.NoError(err)
	suite.Equal(owner2.String(), nft.Owner)
}
