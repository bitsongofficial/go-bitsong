package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestLastAuctionIdGetSet() {
	// get default last auction id
	lastAuctionId := suite.app.MarketplaceKeeper.GetLastAuctionId(suite.ctx)
	suite.Require().Equal(lastAuctionId, uint64(0))

	// set last auction id to new value
	newAuctionId := uint64(2)
	suite.app.MarketplaceKeeper.SetLastAuctionId(suite.ctx, newAuctionId)

	// check last auction id update
	lastAuctionId = suite.app.MarketplaceKeeper.GetLastAuctionId(suite.ctx)
	suite.Require().Equal(lastAuctionId, newAuctionId)
}

func (suite *KeeperTestSuite) TestAuctionGetSet() {
	// get auction by not available id
	_, err := suite.app.MarketplaceKeeper.GetAuctionById(suite.ctx, 0)
	suite.Require().Error(err)

	// get all auctions when not available
	allAuctions := suite.app.MarketplaceKeeper.GetAllAuctions(suite.ctx)
	suite.Require().Len(allAuctions, 0)

	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Second))

	owner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	owner2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	auctions := []types.Auction{
		{ // created auction
			Id:               1,
			Authority:        owner.String(),
			NftId:            1,
			PrizeType:        types.AuctionPrizeType_NftOnlyTransfer,
			Duration:         time.Second,
			BidDenom:         "ubtsg",
			PriceFloor:       1000000,
			InstantSalePrice: 2000000,
			TickSize:         10000,
			State:            types.AuctionState_Created,
			LastBidAmount:    0,
			LastBid:          time.Time{},
			EndedAt:          time.Time{},
			EndAuctionAt:     time.Time{},
			Claimed:          false,
		},
		{ // started auction
			Id:               2,
			Authority:        owner.String(),
			NftId:            2,
			PrizeType:        types.AuctionPrizeType_NftOnlyTransfer,
			Duration:         time.Second,
			BidDenom:         "ubtsg",
			PriceFloor:       1000000,
			InstantSalePrice: 2000000,
			TickSize:         10000,
			State:            types.AuctionState_Started,
			LastBidAmount:    0,
			LastBid:          time.Time{},
			EndedAt:          time.Time{},
			EndAuctionAt:     now.Add(time.Second),
			Claimed:          false,
		},
		{ // bid auction
			Id:               3,
			Authority:        owner.String(),
			NftId:            3,
			PrizeType:        types.AuctionPrizeType_NftOnlyTransfer,
			Duration:         time.Second,
			BidDenom:         "ubtsg",
			PriceFloor:       1000000,
			InstantSalePrice: 2000000,
			TickSize:         10000,
			State:            types.AuctionState_Started,
			LastBidAmount:    1000000,
			LastBid:          now,
			EndedAt:          time.Time{},
			EndAuctionAt:     now.Add(time.Second),
			Claimed:          false,
		},
		{ // ended auction
			Id:               4,
			Authority:        owner2.String(),
			NftId:            4,
			PrizeType:        types.AuctionPrizeType_NftOnlyTransfer,
			Duration:         time.Second,
			BidDenom:         "ubtsg",
			PriceFloor:       1000000,
			InstantSalePrice: 2000000,
			TickSize:         10000,
			State:            types.AuctionState_Ended,
			LastBidAmount:    1000000,
			LastBid:          now,
			EndedAt:          now.Add(time.Second * 2),
			EndAuctionAt:     now.Add(time.Second),
			Claimed:          false,
		},
		{ // claimed auction
			Id:               5,
			Authority:        owner2.String(),
			NftId:            5,
			PrizeType:        types.AuctionPrizeType_NftOnlyTransfer,
			Duration:         time.Second,
			BidDenom:         "ubtsg",
			PriceFloor:       1000000,
			InstantSalePrice: 2000000,
			TickSize:         10000,
			State:            types.AuctionState_Ended,
			LastBidAmount:    1000000,
			LastBid:          now,
			EndedAt:          now.Add(time.Second * 2),
			EndAuctionAt:     now.Add(time.Second),
			Claimed:          true,
		},
	}

	for _, auction := range auctions {
		suite.app.MarketplaceKeeper.SetAuction(suite.ctx, auction)
	}

	for _, auction := range auctions {
		c, err := suite.app.MarketplaceKeeper.GetAuctionById(suite.ctx, auction.Id)
		suite.Require().NoError(err)
		suite.Require().Equal(auction, c)
	}

	allAuctions = suite.app.MarketplaceKeeper.GetAllAuctions(suite.ctx)
	suite.Require().Len(allAuctions, 5)
	suite.Require().Equal(auctions, allAuctions)

	// test GetAuctionsByAuthority
	ownerAuctions := suite.app.MarketplaceKeeper.GetAuctionsByAuthority(suite.ctx, owner)
	suite.Require().Len(ownerAuctions, 3)

	// test DeleteAuction
	suite.app.MarketplaceKeeper.DeleteAuction(suite.ctx, auctions[0])
	allAuctions = suite.app.MarketplaceKeeper.GetAllAuctions(suite.ctx)
	suite.Require().Len(allAuctions, 4)
	ownerAuctions = suite.app.MarketplaceKeeper.GetAuctionsByAuthority(suite.ctx, owner)
	suite.Require().Len(ownerAuctions, 2)

	// auctions to end
	toEndAuctions := suite.app.MarketplaceKeeper.GetAuctionsByAuthority(suite.ctx, owner)
	suite.Require().Len(toEndAuctions, 2)
}

// TODO: test for CreateAuction
// TODO: test for StartAuction
// TODO: test for EndAuction
// TODO: test for SetAuctionAuthority
