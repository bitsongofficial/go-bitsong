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
	creator = sdk.AccAddress(tmhash.SumTruncated([]byte("creator")))
	owner   = sdk.AccAddress(tmhash.SumTruncated([]byte("owner")))
	owner2  = sdk.AccAddress(tmhash.SumTruncated([]byte("owner2")))
	// initAmt  = math.NewIntFromUint64(1000000000)
	// initCoin = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initAmt)}
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	ctx sdk.Context
	// bk     bankkeeper.Keeper
	keeper keeper.Keeper
	app    *simapp.BitsongApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	app := suite.App
	suite.keeper = app.NftKeeper
	// suite.bk = app.BankKeeper
	suite.App = app
	suite.ctx = suite.Ctx

	// init tokens to addr
	/*err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, creator, initCoin)
	suite.NoError(err)*/
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestCreateCollection() {
	testCollection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
	}
	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"

	denom, err := suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.NoError(err)
	suite.Equal(expectedDenom, denom)

	_, err = suite.keeper.CreateCollection(suite.ctx, creator, testCollection)
	suite.Error(err)
}

func (suite *KeeperTestSuite) TestMintNFT() {
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
	suite.Equal(math.NewInt(0), supply)

	nft1 := types.Nft{
		TokenId:     "1",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	nft2 := types.Nft{
		TokenId:     "2",
		Name:        "My First NFT",
		Description: "This is my first NFT",
		Uri:         "ipfs://my-first-nft-metadata.json",
	}

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft1)
	suite.NoError(err)

	supply = suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(math.NewInt(1), supply)

	err = suite.keeper.MintNFT(suite.ctx, collectionDenom, creator, owner, nft2)
	suite.NoError(err)

	supply = suite.keeper.GetSupply(suite.ctx, collectionDenom)
	suite.Equal(math.NewInt(2), supply)
}

func (suite *KeeperTestSuite) TestSendNFT() {
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

	res, err := suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 1)
	suite.Equal(nft1.TokenId, res.Nfts[0].TokenId)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 0)

	err = suite.keeper.SendNft(suite.ctx, owner, owner2, collectionDenom, "1")
	suite.NoError(err)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner2.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 1)
	suite.Equal(nft1.TokenId, res.Nfts[0].TokenId)

	res, err = suite.keeper.AllNftsByOwner(suite.ctx, &types.QueryAllNftsByOwnerRequest{
		Owner: owner.String(),
	})
	suite.NoError(err)
	suite.Len(res.Nfts, 0)

	nft, err := suite.keeper.NFTs.Get(suite.ctx, collections.Join(collectionDenom, "1"))
	suite.NoError(err)
	suite.Equal(owner2.String(), nft.Owner)
}
