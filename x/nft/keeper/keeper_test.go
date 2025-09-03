package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	simapp "github.com/bitsongofficial/go-bitsong/app"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/suite"
)

var (
	creator  = sdk.AccAddress(tmhash.SumTruncated([]byte("creator")))
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("owner")))
	initAmt  = math.NewIntFromUint64(1000000000)
	initCoin = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initAmt)}
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	ctx    sdk.Context
	bk     bankkeeper.Keeper
	keeper keeper.Keeper
	app    *simapp.BitsongApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	app := suite.App
	suite.keeper = app.NftKeeper
	suite.bk = app.BankKeeper
	suite.App = app
	suite.ctx = suite.Ctx

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, fantokentypes.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, fantokentypes.ModuleName, creator, initCoin)
	suite.NoError(err)
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

/*
type MintNFTTestCase struct {
	name                       string // A descriptive name for the test case
	collection                 types.Collection
	nftToMint                  types.Nft
	minter                     sdk.AccAddress
	owner                      sdk.AccAddress
	expectErr                  bool  // Do we expect an error during minting?
	expectedSupply             int64 // Expected supply of the collection after this mint
	expectedBalanceForNewNFT   int64 // Expected balance of the *specific* new NFT
	expectedTotalOwnerBalances int   // Expected total number of different assets the owner has
}

func (suite *KeeperTestSuite) TestMintNFT_Advanced() {
	collection := types.Collection{
		Name:        "My NFT Collection",
		Symbol:      "MYNFT",
		Description: "My NFT Collection Description",
		Uri:         "ipfs://my-nft-collection-metadata.json",
		Minter:      creator.String(),
	}

	collectionDenom, err := suite.keeper.CreateCollection(suite.ctx, creator, collection)
	suite.Require().NoError(err, "initial collection creation should succeed")
	fmt.Println("collectionDenom:", collectionDenom)

	// Define the test cases
	testCases := []MintNFTTestCase{
		{
			name:       "Successful first mint",
			collection: collection,
			nftToMint: types.Nft{
				Name:        "My First NFT",
				Description: "This is my first NFT",
				Uri:         "ipfs://my-first-nft-metadata.json",
			},
			minter:                     creator,
			owner:                      owner,
			expectErr:                  false,
			expectedSupply:             1,
			expectedBalanceForNewNFT:   1,
			expectedTotalOwnerBalances: 1, // Owner should now have 1 asset.
		},
		{
			name:       "Successful second mint",
			collection: collection,
			nftToMint: types.Nft{
				Name:        "My Second NFT",
				Description: "This is my second NFT",
				Uri:         "ipfs://my-second-nft-metadata.json",
			},
			minter:                     creator,
			owner:                      owner,
			expectErr:                  false,
			expectedSupply:             2, // Total supply of the collection is now 2.
			expectedBalanceForNewNFT:   1,
			expectedTotalOwnerBalances: 2, // Owner now has two distinct NFTs.
		},
		{
			name:       "Unauthorized minter should fail",
			collection: collection,
			nftToMint: types.Nft{
				Name:        "Unauthorized NFT",
				Description: "This NFT should not be minted",
				Uri:         "ipfs://unauthorized-nft-metadata.json",
			},
			minter:                     owner, // 'owner' is not the authorized minter
			owner:                      owner,
			expectErr:                  true,
			expectedSupply:             2, // Supply should NOT increase.
			expectedBalanceForNewNFT:   0, // Not applicable, but setting to 0 for clarity.
			expectedTotalOwnerBalances: 2, // Owner's total assets should remain unchanged.
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			nftDenom, err := suite.keeper.MintNFT(suite.ctx, collectionDenom, tc.minter, tc.owner, tc.nftToMint)

			if tc.expectErr {
				suite.Require().Error(err, "should have returned an error")
			} else {
				suite.Require().NoError(err, "should not have returned an error")
				fmt.Println("nftDenom for '"+tc.name+"':", nftDenom)

				// Check the balance of the newly minted NFT for the owner
				resp, err := suite.bk.Balance(suite.ctx, &types2.QueryBalanceRequest{
					Address: tc.owner.String(),
					Denom:   nftDenom,
				})
				suite.NoError(err)
				suite.Equal(tc.expectedBalanceForNewNFT, resp.Balance.Amount.Int64(), "owner's balance for the new NFT should match expected")
			}

			// Check the supply of the collection
			supply := suite.keeper.GetSupply(suite.ctx, collectionDenom)
			suite.Equal(math.NewInt(tc.expectedSupply), supply, "collection supply should match expected")

			// Check the owner's total number of different assets
			balances, err := suite.bk.AllBalances(suite.ctx, &types2.QueryAllBalancesRequest{
				Address: tc.owner.String(),
			})
			suite.NoError(err, "querying all balances should not produce an error")
			suite.Equal(tc.expectedTotalOwnerBalances, len(balances.Balances), "owner's total number of assets should match expected")
		})
	}
}
*/
