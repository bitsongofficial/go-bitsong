package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
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

func (suite *KeeperTestSuite) TestPlaceBid() {
	suite.ctx = suite.ctx.WithBlockTime(time.Now().UTC())
	owner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	coins := sdk.Coins{sdk.NewInt64Coin("ubtsg", 1000000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, bidder, coins)

	tests := []struct {
		testCase         string
		state            types.AuctionState
		isLastBidder     bool
		lastBidAmount    uint64
		bidToken         string
		newBidAmount     uint64
		instantSalePrice uint64
		auctionId        uint64
		expectPass       bool
	}{
		{
			"Not existing auction id",
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1500,
			10000,
			0,
			false,
		},
		{
			"bid on not active auction",
			types.AuctionState_Created,
			false,
			1000,
			"ubtsg",
			1500,
			10000,
			1,
			false,
		},
		{
			"invalid bid token",
			types.AuctionState_Ended,
			false,
			1000,
			"randtoken",
			1500,
			10000,
			1,
			false,
		},
		{
			"bid with low amount check",
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			100,
			10000,
			1,
			false,
		},
		{
			"bid by winner bidder check",
			types.AuctionState_Started,
			true,
			1000,
			"ubtsg",
			0,
			10000,
			1,
			false,
		},
		{
			"successful bid lower than instant sale price",
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1500,
			10000,
			1,
			true,
		},
		{
			"successful bid higher than instant sale price",
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			11000,
			10000,
			1,
			true,
		},
	}

	for _, tc := range tests {
		// module account
		moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

		// set nft with ownership
		nft := nfttypes.NFT{
			Id:         1,
			Owner:      moduleAddr.String(),
			MetadataId: 1,
		}
		suite.app.NFTKeeper.SetNFT(suite.ctx, nft)

		// set metadata with ownership
		metadata := nfttypes.Metadata{
			Id:              1,
			UpdateAuthority: moduleAddr.String(),
		}
		suite.app.NFTKeeper.SetMetadata(suite.ctx, metadata)

		// set auction with ownership
		auction := types.Auction{
			Id:            1,
			Authority:     owner.String(),
			NftId:         1,
			Duration:      time.Second,
			PrizeType:     types.AuctionPrizeType_NftOnlyTransfer,
			State:         tc.state,
			LastBidAmount: tc.lastBidAmount,
			BidDenom:      "ubtsg",
		}
		suite.app.MarketplaceKeeper.SetAuction(suite.ctx, auction)

		if tc.isLastBidder {
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.lastBidAmount,
			})
		} else {
			suite.app.MarketplaceKeeper.DeleteBid(suite.ctx, types.Bid{
				Bidder:    bidder.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.lastBidAmount,
			})
		}

		oldBidderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bidder, "ubtsg")
		oldModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, "ubtsg")

		// execute SetAuctionAuthority
		msg := types.NewMsgPlaceBid(bidder, tc.auctionId, sdk.NewInt64Coin(tc.bidToken, int64(tc.newBidAmount)))
		err := suite.app.MarketplaceKeeper.PlaceBid(suite.ctx, msg)

		// check error exists on the execution
		if tc.expectPass {
			suite.Require().NoError(err)

			// check bid object is added
			bid, err := suite.app.MarketplaceKeeper.GetBid(suite.ctx, tc.auctionId, bidder)
			suite.Require().NoError(err)
			suite.Require().Equal(bid.AuctionId, uint64(1))
			suite.Require().Equal(bid.Amount, tc.newBidAmount)

			newBidderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bidder, "ubtsg")
			newModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, "ubtsg")

			// check balance has been reduced from end user
			suite.Require().Equal(oldBidderBalance.Amount, newBidderBalance.Amount.Add(sdk.NewInt(int64(tc.newBidAmount))))

			// check balance has been increased on module account
			suite.Require().Equal(newModuleBalance.Amount, oldModuleBalance.Amount.Add(sdk.NewInt(int64(tc.newBidAmount))))

			// check updated auction object with lastBid and lastBidAmount
			auction, err := suite.app.MarketplaceKeeper.GetAuctionById(suite.ctx, tc.auctionId)
			suite.Require().NoError(err)
			suite.Require().Equal(auction.LastBid, suite.ctx.BlockTime())
			suite.Require().Equal(auction.LastBidAmount, tc.newBidAmount)

			// check bidder metadata has been set correctly
			biddermeta, err := suite.app.MarketplaceKeeper.GetBidderMetadata(suite.ctx, bidder)
			suite.Require().NoError(err)
			suite.Require().Equal(biddermeta, types.BidderMetadata{
				Bidder:           msg.Sender,
				LastAuctionId:    msg.AuctionId,
				LastBid:          msg.Amount.Amount.Uint64(),
				LastBidTimestamp: suite.ctx.BlockTime(),
				LastBidCancelled: false,
			})

			// check auction end when it's bigger than instant sale price
			if auction.InstantSalePrice <= tc.newBidAmount {
				suite.Require().Equal(auction.State, types.AuctionState_Ended)
			}
		} else {
			suite.Require().Error(err)
		}
	}
}

// TODO: add test for CancelBid
// TODO: add test for ClaimBid
