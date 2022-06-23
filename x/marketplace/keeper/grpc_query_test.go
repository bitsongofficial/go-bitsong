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
			NftId:            "1",
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
			Claimed:          0,
		},
		{ // started auction
			Id:               2,
			Authority:        owner.String(),
			NftId:            "2",
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
			Claimed:          0,
		},
		{ // bid auction
			Id:               3,
			Authority:        owner.String(),
			NftId:            "3",
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
			Claimed:          0,
		},
		{ // ended auction
			Id:               4,
			Authority:        owner2.String(),
			NftId:            "4",
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
			Claimed:          0,
		},
		{ // claimed auction
			Id:               5,
			Authority:        owner2.String(),
			NftId:            "5",
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
			Claimed:          1,
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
			NftId:            "1",
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
			Claimed:          0,
		},
		{ // started auction
			Id:               2,
			Authority:        owner.String(),
			NftId:            "2",
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
			Claimed:          0,
		},
		{ // bid auction
			Id:               3,
			Authority:        owner.String(),
			NftId:            "3",
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
			Claimed:          0,
		},
		{ // ended auction
			Id:               4,
			Authority:        owner2.String(),
			NftId:            "4",
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
			Claimed:          0,
		},
		{ // claimed auction
			Id:               5,
			Authority:        owner2.String(),
			NftId:            "5",
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
			Claimed:          1,
		},
	}

	for _, auction := range auctions {
		suite.app.MarketplaceKeeper.SetAuction(suite.ctx, auction)
	}

	for _, auction := range auctions {
		c, err := suite.app.MarketplaceKeeper.Auction(sdk.WrapSDKContext(suite.ctx), &types.QueryAuctionRequest{
			Id: auction.Id,
		})
		suite.Require().NoError(err)
		suite.Require().Equal(auction, c.Auction)
	}
}

func (suite *KeeperTestSuite) TestGRPCBidsByAuction() {
	bidder1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Second))

	// get not available bid
	_, err := suite.app.MarketplaceKeeper.GetBid(suite.ctx, 1, bidder1)
	suite.Require().Error(err)

	// get all bids when not available
	allBids := suite.app.MarketplaceKeeper.GetAllBids(suite.ctx)
	suite.Require().Len(allBids, 0)

	bids := []types.Bid{
		{
			Bidder:    bidder1.String(),
			AuctionId: 1,
			Amount:    1000000,
		},
		{
			Bidder:    bidder2.String(),
			AuctionId: 1,
			Amount:    1200000,
		},
		{
			Bidder:    bidder1.String(),
			AuctionId: 2,
			Amount:    1000000,
		},
	}

	for _, bid := range bids {
		suite.app.MarketplaceKeeper.SetBid(suite.ctx, bid)
	}

	// test GetBidsAuction
	resp, err := suite.app.MarketplaceKeeper.BidsByAuction(sdk.WrapSDKContext(suite.ctx), &types.QueryBidsByAuctionRequest{
		Id: 1,
	})
	suite.Require().Len(resp.Bids, 2)
}

func (suite *KeeperTestSuite) TestGRPCBidsByBidder() {
	bidder1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Second))

	bids := []types.Bid{
		{
			Bidder:    bidder1.String(),
			AuctionId: 1,
			Amount:    1000000,
		},
		{
			Bidder:    bidder2.String(),
			AuctionId: 1,
			Amount:    1200000,
		},
		{
			Bidder:    bidder1.String(),
			AuctionId: 2,
			Amount:    1000000,
		},
	}

	for _, bid := range bids {
		suite.app.MarketplaceKeeper.SetBid(suite.ctx, bid)
	}

	// test GetBidsByBidder
	resp, err := suite.app.MarketplaceKeeper.BidsByBidder(sdk.WrapSDKContext(suite.ctx), &types.QueryBidsByBidderRequest{
		Bidder: bidder1.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Bids, 2)
}

func (suite *KeeperTestSuite) TestGRPCBidderMetadata() {

	bidder1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Second))

	bidderdata := []types.BidderMetadata{
		{
			Bidder:           bidder1.String(),
			LastAuctionId:    1,
			LastBid:          5000,
			LastBidTimestamp: now.UTC(),
			LastBidCancelled: false,
		},
		{
			Bidder:           bidder2.String(),
			LastAuctionId:    1,
			LastBid:          10000,
			LastBidTimestamp: now.UTC(),
			LastBidCancelled: false,
		},
	}

	for _, bm := range bidderdata {
		suite.app.MarketplaceKeeper.SetBidderMetadata(suite.ctx, bm)
	}

	for _, bm := range bidderdata {
		bidder, err := sdk.AccAddressFromBech32(bm.Bidder)
		suite.Require().NoError(err)
		b, err := suite.app.MarketplaceKeeper.BidderMetadata(sdk.WrapSDKContext(suite.ctx), &types.QueryBidderMetadataRequest{
			Bidder: bidder.String(),
		})
		suite.Require().NoError(err)
		suite.Require().Equal(b.BidderMetadata, bm)
	}
}
