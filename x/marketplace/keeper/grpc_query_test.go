package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestGRPCAuctions() {
	// get all auctions when not available
	resp, err := suite.app.MarketplaceKeeper.Auctions(sdk.WrapSDKContext(suite.ctx), &types.QueryAuctionsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Auctions, 0)

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

	// all auctions
	resp, err = suite.app.MarketplaceKeeper.Auctions(sdk.WrapSDKContext(suite.ctx), &types.QueryAuctionsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Auctions, 5)
	suite.Require().Equal(auctions, resp.Auctions)

	// auctions by authority
	resp, err = suite.app.MarketplaceKeeper.Auctions(sdk.WrapSDKContext(suite.ctx), &types.QueryAuctionsRequest{
		Authority: owner.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Auctions, 3)

	// auctions by state
	resp, err = suite.app.MarketplaceKeeper.Auctions(sdk.WrapSDKContext(suite.ctx), &types.QueryAuctionsRequest{
		State: types.AuctionState_Started,
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Auctions, 2)
}

func (suite *KeeperTestSuite) TestGRPCAuction() {
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
		c, err := suite.app.MarketplaceKeeper.Auction(sdk.WrapSDKContext(suite.ctx), &types.QueryAuctionRequest{
			Id: auction.NftId,
		})
		suite.Require().NoError(err)
		suite.Require().Equal(auction, c.Auction)
	}
}

// TODO: add test for BidsByAuction
// TODO: add test for BidsByBidder
// TODO: add test for BidderMetadata
