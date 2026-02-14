package keeper_test

import (
	"strings"
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

	authority1 = sdk.AccAddress(tmhash.SumTruncated([]byte("authority1")))
	authority2 = sdk.AccAddress(tmhash.SumTruncated([]byte("authority2")))

	testCollection1 = types.Collection{
		Name:   "My NFT Collection",
		Symbol: "MYNFT",
		Uri:    "ipfs://my-nft-collection-metadata.json",
		Minter: minter1.String(),
	}
	expectedDenom1 = "nft9436DDD23FB751AEA7BC6C767F20F943DD735E06"

	testCollection2 = types.Collection{
		Name:   "My Second NFT Collection",
		Symbol: "MYNFT2",
		Uri:    "ipfs://my-nft-collection-metadata.json",
		Minter: minter1.String(),
	}
	expectedDenom2 = "nft6C5B4EC9EA22932F217B0A0CDCA3A987B8271CD0"

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

	msgServer types.MsgServer

	app *simapp.BitsongApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	app := suite.App
	suite.keeper = app.NftKeeper
	suite.App = app
	suite.ctx = suite.Ctx

	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestCreateCollection() {
	denom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		"",
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom1, denom)

	_, err = suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.Error(err)

	denom, err = suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
		testCollection2.Symbol,
		testCollection2.Name,
		testCollection2.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom2, denom)
}

func (suite *KeeperTestSuite) TestSetCollectionName() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom1, collectionDenom)

	err = suite.keeper.SetCollectionName(suite.ctx, authority1, collectionDenom, "New Collection Name")
	suite.NoError(err)

	collection, err := suite.keeper.GetCollection(suite.ctx, collectionDenom)
	suite.NoError(err)
	suite.Equal("New Collection Name", collection.Name)

	err = suite.keeper.SetCollectionName(suite.ctx, creator2, collectionDenom, "Another Collection Name")
	suite.Error(err)

	err = suite.keeper.SetCollectionName(suite.ctx, creator1, collectionDenom, strings.Repeat("a", 65))
	suite.Error(err)

	// test setting name on a collection without authority
	collectionDenom2, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter2.String(),
		"",
		testCollection2.Symbol,
		testCollection2.Name,
		testCollection2.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom2, collectionDenom2)

	err = suite.keeper.SetCollectionName(suite.ctx, creator1, collectionDenom2, "New Collection Name")
	suite.Error(err)
	suite.Contains(err.Error(), "only the collection authority can change the name")

	// create a collection without minter and authority
	_, err = suite.keeper.CreateCollection(
		suite.ctx,
		creator2.String(),
		"",
		"",
		"COLL3",
		"Collection 3",
		"ipfs://collection-3-metadata.json",
	)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestSetCollectionUri() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom1, collectionDenom)

	err = suite.keeper.SetCollectionUri(suite.ctx, authority1, collectionDenom, "ipfs://new-uri.json")
	suite.NoError(err)

	collection, err := suite.keeper.GetCollection(suite.ctx, collectionDenom)
	suite.NoError(err)
	suite.Equal("ipfs://new-uri.json", collection.Uri)

	err = suite.keeper.SetCollectionUri(suite.ctx, creator2, collectionDenom, "ipfs://another-uri.json")
	suite.Error(err)

	err = suite.keeper.SetCollectionUri(suite.ctx, authority1, collectionDenom, strings.Repeat("a", 165))
	suite.Error(err)

	// test setting uri on a collection without authority
	collectionDenom2, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter2.String(),
		"",
		testCollection2.Symbol,
		testCollection2.Name,
		testCollection2.Uri,
	)
	suite.NoError(err)
	suite.Equal(expectedDenom2, collectionDenom2)

	err = suite.keeper.SetCollectionUri(suite.ctx, creator1, collectionDenom2, "ipfs://new-uri.json")
	suite.Error(err)
	suite.Contains(err.Error(), "only the collection authority can change the uri")
}

func (suite *KeeperTestSuite) TestMintNFT() {
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

	// collection with no minter
	collectionDenom2, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator2.String(),
		"",
		"",
		testCollection2.Symbol,
		testCollection2.Name,
		testCollection2.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.MintNFT(
		suite.ctx,
		creator2,
		owner2,
		collectionDenom2,
		testNft3.TokenId,
		testNft3.Name,
		testNft3.Uri,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "minting disabled for this collection")
}

func (suite *KeeperTestSuite) TestSetMinter() {
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

	err = suite.keeper.SetMinter(suite.ctx, minter1, &minter2, collectionDenom)
	suite.NoError(err)

	minter, err := suite.keeper.GetMinter(suite.ctx, collectionDenom)
	suite.NoError(err)
	suite.Equal(minter2.String(), minter.String())

	err = suite.keeper.SetMinter(suite.ctx, minter1, &minter1, collectionDenom)
	suite.Error(err)

	err = suite.keeper.SetMinter(suite.ctx, minter2, nil, collectionDenom)
	suite.NoError(err)

	_, err = suite.keeper.GetMinter(suite.ctx, collectionDenom)
	suite.Error(err)
	suite.Contains(err.Error(), "minting disabled for this collection")

	err = suite.keeper.SetMinter(suite.ctx, minter2, &minter1, collectionDenom)
	suite.Error(err)
	suite.Contains(err.Error(), "minting disabled for this collection")
}

func (suite *KeeperTestSuite) TestSetAuthority() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
		testCollection1.Symbol,
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.SetAuthority(suite.ctx, authority1, &authority2, collectionDenom)
	suite.NoError(err)

	collection, err := suite.keeper.GetCollection(suite.ctx, collectionDenom)
	suite.NoError(err)
	suite.Equal(authority2.String(), collection.Authority)

	err = suite.keeper.SetAuthority(suite.ctx, authority1, &creator2, collectionDenom)
	suite.Error(err)
	suite.Contains(err.Error(), "only the current authority can change the authority")

	err = suite.keeper.SetAuthority(suite.ctx, authority2, nil, collectionDenom)
	suite.NoError(err)

	collection, err = suite.keeper.GetCollection(suite.ctx, collectionDenom)
	suite.NoError(err)
	suite.Equal("", collection.Authority)

	err = suite.keeper.SetAuthority(suite.ctx, authority2, &creator2, collectionDenom)
	suite.Error(err)
	suite.Contains(err.Error(), "only the current authority can change the authority")
}

func (suite *KeeperTestSuite) TestSendNFT() {
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

	err = suite.keeper.SendNFT(suite.ctx, owner1, owner2, collectionDenom, "1")
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

func (suite *KeeperTestSuite) TestSetNFTName() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
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

	err = suite.keeper.SetNFTName(suite.ctx, authority1, collectionDenom, testNft1.TokenId, "New NFT Name")
	suite.NoError(err)

	nft, err := suite.keeper.NFTs.Get(suite.ctx, collections.Join(collectionDenom, testNft1.TokenId))
	suite.NoError(err)
	suite.Equal("New NFT Name", nft.Name)

	err = suite.keeper.SetNFTName(suite.ctx, owner2, collectionDenom, testNft1.TokenId, "Another NFT Name")
	suite.Error(err)
	suite.Contains(err.Error(), "only the collection authority can set NFT name")

	err = suite.keeper.SetNFTName(suite.ctx, authority1, collectionDenom, testNft1.TokenId, strings.Repeat("a", 165))
	suite.Error(err)
	suite.Contains(err.Error(), "name length exceeds maximum")

	// test setting name on a collection without authority
	collectionDenom2, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator2.String(),
		minter2.String(),
		"",
		testCollection2.Symbol,
		testCollection2.Name,
		testCollection2.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter2,
		owner2,
		collectionDenom2,
		testNft2.TokenId,
		testNft2.Name,
		testNft2.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.SetNFTName(suite.ctx, owner2, collectionDenom2, testNft2.TokenId, "New NFT Name")
	suite.Error(err)
	suite.Contains(err.Error(), "no authority, cannot set NFT name")

	// create a collection with authority
	collectionDenom3, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
		testCollection1.Symbol+"3",
		testCollection1.Name,
		testCollection1.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter1,
		owner1,
		collectionDenom3,
		testNft3.TokenId,
		testNft3.Name,
		testNft3.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.SetNFTName(suite.ctx, owner1, collectionDenom3, testNft3.TokenId, "New NFT Name")
	suite.Error(err)
	suite.Contains(err.Error(), "only the collection authority can set NFT name")

	err = suite.keeper.SetNFTName(suite.ctx, authority1, collectionDenom3, testNft3.TokenId, "New NFT Name")
	suite.NoError(err)

	nft, err = suite.keeper.NFTs.Get(suite.ctx, collections.Join(collectionDenom3, testNft3.TokenId))
	suite.NoError(err)
	suite.Equal("New NFT Name", nft.Name)
}

func (suite *KeeperTestSuite) TestSetNFTUri() {
	collectionDenom, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator1.String(),
		minter1.String(),
		authority1.String(),
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

	err = suite.keeper.SetNFTUri(suite.ctx, authority1, collectionDenom, testNft1.TokenId, "ipfs://new-uri.json")
	suite.NoError(err)

	nft, err := suite.keeper.NFTs.Get(suite.ctx, collections.Join(collectionDenom, testNft1.TokenId))
	suite.NoError(err)
	suite.Equal("ipfs://new-uri.json", nft.Uri)

	err = suite.keeper.SetNFTUri(suite.ctx, owner2, collectionDenom, testNft1.TokenId, "ipfs://another-uri.json")
	suite.Error(err)
	suite.Contains(err.Error(), "only the collection authority can set NFT uri")

	err = suite.keeper.SetNFTUri(suite.ctx, authority1, collectionDenom, testNft1.TokenId, strings.Repeat("a", 165))
	suite.Error(err)
	suite.Contains(err.Error(), "URI length exceeds maximum")

	// test setting uri on a collection without authority
	collectionDenom2, err := suite.keeper.CreateCollection(
		suite.ctx,
		creator2.String(),
		minter2.String(),
		"",
		testCollection2.Symbol,
		testCollection2.Name,
		testCollection2.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.MintNFT(
		suite.ctx,
		minter2,
		owner2,
		collectionDenom2,
		testNft2.TokenId,
		testNft2.Name,
		testNft2.Uri,
	)
	suite.NoError(err)

	err = suite.keeper.SetNFTUri(suite.ctx, owner2, collectionDenom2, testNft2.TokenId, "ipfs://new-uri.json")
	suite.Error(err)
	suite.Contains(err.Error(), "no authority, cannot set NFT uri")
}
