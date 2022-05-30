package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestBidGetSet() {
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

	for _, bid := range bids {
		bidder, err := sdk.AccAddressFromBech32(bid.Bidder)
		suite.Require().NoError(err)
		b, err := suite.app.MarketplaceKeeper.GetBid(suite.ctx, bid.AuctionId, bidder)
		suite.Require().NoError(err)
		suite.Require().Equal(bid, b)
	}

	allBids = suite.app.MarketplaceKeeper.GetAllBids(suite.ctx)
	suite.Require().Len(allBids, 3)

	// test GetBidsByBidder
	bidder1Bids := suite.app.MarketplaceKeeper.GetBidsByBidder(suite.ctx, bidder1)
	suite.Require().Len(bidder1Bids, 2)

	// test GetBidsByBidder
	auction1Bids := suite.app.MarketplaceKeeper.GetBidsByAuction(suite.ctx, 1)
	suite.Require().Len(auction1Bids, 2)

	// test DeleteBid
	suite.app.MarketplaceKeeper.DeleteBid(suite.ctx, bids[0])
	allBids = suite.app.MarketplaceKeeper.GetAllBids(suite.ctx)
	suite.Require().Len(allBids, 2)
	bidder1Bids = suite.app.MarketplaceKeeper.GetBidsByBidder(suite.ctx, bidder1)
	suite.Require().Len(bidder1Bids, 1)
	auction1Bids = suite.app.MarketplaceKeeper.GetBidsByAuction(suite.ctx, 1)
	suite.Require().Len(auction1Bids, 1)
}

func (suite *KeeperTestSuite) TestBidderMetadataGetSet() {

	bidder1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Second))

	// get not available bidder metadata
	_, err := suite.app.MarketplaceKeeper.GetBidderMetadata(suite.ctx, bidder1)
	suite.Require().Error(err)

	// get all bidder metadata when not available
	allBidderData := suite.app.MarketplaceKeeper.GetAllBidderMetadata(suite.ctx)
	suite.Require().Len(allBidderData, 0)

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
		b, err := suite.app.MarketplaceKeeper.GetBidderMetadata(suite.ctx, bidder)
		suite.Require().NoError(err)
		suite.Require().Equal(b, bm)
	}

	allBidderData = suite.app.MarketplaceKeeper.GetAllBidderMetadata(suite.ctx)
	suite.Require().Len(allBidderData, 2)
}

// TODO: add test for PlaceBid
// TODO: add test for CancelBid
// TODO: add test for ClaimBid
