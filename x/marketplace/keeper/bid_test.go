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
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	coins := sdk.Coins{sdk.NewInt64Coin("ubtsg", 1000000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, bidder, coins)

	tests := []struct {
		testCase         string
		auctionType      types.AuctionPrizeType
		state            types.AuctionState
		isLastBidder     bool
		lastBidAmount    uint64
		bidToken         string
		newBidAmount     uint64
		instantSalePrice uint64
		editionLimit     uint64
		auctionId        uint64
		expectPass       bool
	}{
		{
			"Not existing auction id",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1500,
			10000,
			0,
			0,
			false,
		},
		{
			"bid on not active auction",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Created,
			false,
			1000,
			"ubtsg",
			1500,
			10000,
			0,
			1,
			false,
		},
		{
			"invalid bid token",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Ended,
			false,
			1000,
			"randtoken",
			1500,
			10000,
			0,
			1,
			false,
		},
		{
			"bid with low amount check",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			100,
			10000,
			0,
			1,
			false,
		},
		{
			"bid by winner bidder check",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			true,
			1000,
			"ubtsg",
			0,
			10000,
			0,
			1,
			false,
		},
		{
			"successful bid lower than instant sale price",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1500,
			10000,
			0,
			1,
			true,
		},
		{
			"successful bid higher than instant sale price",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			11000,
			10000,
			0,
			1,
			true,
		},
		{
			"not successful bid for exceeding edition limit",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1000,
			10000,
			1,
			1,
			false,
		},
		{
			"successful bid on limited edition",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1000,
			10000,
			2,
			1,
			true,
		},
		{
			"successful bid on limited edition",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1000,
			10000,
			2,
			1,
			true,
		},
		{
			"not successful bid on open edition for floor price",
			types.AuctionPrizeType_OpenEditionPrints,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			100,
			10000,
			0,
			1,
			false,
		},
		{
			"not successful bid on limited edition for floor price",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			100,
			10000,
			0,
			1,
			false,
		},
		{
			"successful bid on open edition",
			types.AuctionPrizeType_OpenEditionPrints,
			types.AuctionState_Started,
			false,
			1000,
			"ubtsg",
			1000,
			10000,
			0,
			1,
			true,
		},
	}

	for _, tc := range tests {
		// module account
		moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

		// set nft with ownership
		nft := nfttypes.NFT{
			CollId:     1,
			MetadataId: 1,
			Seq:        0,
			Owner:      moduleAddr.String(),
		}
		suite.app.NFTKeeper.SetNFT(suite.ctx, nft)

		// set metadata with ownership
		metadata := nfttypes.Metadata{
			Id:                1,
			MetadataAuthority: moduleAddr.String(),
			MintAuthority:     moduleAddr.String(),
		}
		suite.app.NFTKeeper.SetMetadata(suite.ctx, metadata)

		// set auction with ownership
		auction := types.Auction{
			Id:            1,
			Authority:     owner.String(),
			NftId:         nft.Id(),
			Duration:      time.Second,
			PrizeType:     tc.auctionType,
			State:         tc.state,
			LastBidAmount: tc.lastBidAmount,
			BidDenom:      "ubtsg",
			PriceFloor:    500,
			EditionLimit:  tc.editionLimit,
		}
		suite.app.MarketplaceKeeper.SetAuction(suite.ctx, auction)

		bids := suite.app.MarketplaceKeeper.GetAllBids(suite.ctx)
		for _, bid := range bids {
			suite.app.MarketplaceKeeper.DeleteBid(suite.ctx, bid)
		}
		if tc.isLastBidder {
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.lastBidAmount,
			})
		} else {
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder2.String(),
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
			switch tc.auctionType {
			case types.AuctionPrizeType_NftOnlyTransfer:
				fallthrough
			case types.AuctionPrizeType_FullRightsTransfer:
				if auction.InstantSalePrice <= tc.newBidAmount {
					suite.Require().Equal(auction.State, types.AuctionState_Ended)
				}
			}
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestCancelBid() {
	suite.ctx = suite.ctx.WithBlockTime(time.Now().UTC())
	owner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	coins := sdk.Coins{sdk.NewInt64Coin("ubtsg", 1000000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.ModuleName, coins)

	tests := []struct {
		testCase      string
		auctionType   types.AuctionPrizeType
		state         types.AuctionState
		lastBidAmount uint64
		anotherBid    uint64
		bidAmount     uint64
		editionLimit  uint64
		auctionId     uint64
		expectPass    bool
	}{
		{
			"Not existing auction id",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			1500,
			0,
			1000,
			0,
			0,
			false,
		},
		{
			"cancelling not existing bid",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Created,
			1500,
			0,
			0,
			0,
			1,
			false,
		},
		{
			"try to cancel winner bid",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Ended,
			1000,
			0,
			1000,
			0,
			1,
			false,
		},
		{
			"successful bid cancel",
			types.AuctionPrizeType_NftOnlyTransfer,
			types.AuctionState_Started,
			1000,
			0,
			100,
			0,
			1,
			true,
		},
		{
			"not successful cancel for open edition",
			types.AuctionPrizeType_OpenEditionPrints,
			types.AuctionState_Started,
			1000,
			0,
			1000,
			0,
			1,
			false,
		},
		{
			"not successful cancel for limited edition",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			1000,
			0,
			1000,
			1,
			1,
			false,
		},
		{
			"not successful cancel for limited edition",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			1000,
			0,
			1000,
			1,
			1,
			false,
		},
		{
			"successful cancel of limited edition",
			types.AuctionPrizeType_LimitedEditionPrints,
			types.AuctionState_Started,
			1000,
			1100,
			1000,
			1,
			1,
			true,
		},
	}

	for _, tc := range tests {
		// module account
		moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

		// set nft with ownership
		nft := nfttypes.NFT{
			CollId:     1,
			MetadataId: 1,
			Seq:        0,
			Owner:      moduleAddr.String(),
		}
		suite.app.NFTKeeper.SetNFT(suite.ctx, nft)

		// set metadata with ownership
		metadata := nfttypes.Metadata{
			Id:                1,
			MetadataAuthority: moduleAddr.String(),
			MintAuthority:     moduleAddr.String(),
		}
		suite.app.NFTKeeper.SetMetadata(suite.ctx, metadata)

		// set auction with ownership
		auction := types.Auction{
			Id:            1,
			Authority:     owner.String(),
			NftId:         nft.Id(),
			Duration:      time.Second,
			PrizeType:     tc.auctionType,
			State:         tc.state,
			LastBidAmount: tc.lastBidAmount,
			BidDenom:      "ubtsg",
			PriceFloor:    500,
			EditionLimit:  tc.editionLimit,
		}
		suite.app.MarketplaceKeeper.SetAuction(suite.ctx, auction)

		bids := suite.app.MarketplaceKeeper.GetAllBids(suite.ctx)
		for _, bid := range bids {
			suite.app.MarketplaceKeeper.DeleteBid(suite.ctx, bid)
		}

		if tc.anotherBid > 0 {
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder2.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.anotherBid,
			})
		}
		if tc.bidAmount > 0 {
			suite.app.MarketplaceKeeper.SetBidderMetadata(suite.ctx, types.BidderMetadata{
				Bidder:           bidder.String(),
				LastAuctionId:    tc.auctionId,
				LastBid:          tc.bidAmount,
				LastBidTimestamp: suite.ctx.BlockTime(),
				LastBidCancelled: false,
			})
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.bidAmount,
			})
		}

		oldBidderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bidder, "ubtsg")
		oldModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, "ubtsg")

		// execute SetAuctionAuthority
		msg := types.NewMsgCancelBid(bidder, tc.auctionId)
		err := suite.app.MarketplaceKeeper.CancelBid(suite.ctx, msg)

		// check error exists on the execution
		if tc.expectPass {
			suite.Require().NoError(err)

			// check bid object is deleted
			_, err := suite.app.MarketplaceKeeper.GetBid(suite.ctx, tc.auctionId, bidder)
			suite.Require().Error(err)

			newBidderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bidder, "ubtsg")
			newModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, "ubtsg")

			// check balance has been reduced from end user
			suite.Require().Equal(newBidderBalance.Amount, oldBidderBalance.Amount.Add(sdk.NewInt(int64(tc.bidAmount))))

			// check balance has been increased on module account
			suite.Require().Equal(oldModuleBalance.Amount, newModuleBalance.Amount.Add(sdk.NewInt(int64(tc.bidAmount))))

			// check bidder metadata is updated for LastBidCancelled true
			biddermeta, err := suite.app.MarketplaceKeeper.GetBidderMetadata(suite.ctx, bidder)
			suite.Require().NoError(err)
			suite.Require().Equal(biddermeta.LastBidCancelled, true)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestClaimBid() {
	suite.ctx = suite.ctx.WithBlockTime(time.Now().UTC())
	owner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	bidder2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	coins := sdk.Coins{sdk.NewInt64Coin("ubtsg", 1000000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.ModuleName, coins)

	tests := []struct {
		testCase       string
		state          types.AuctionState
		prizeType      types.AuctionPrizeType
		presaleHappend bool
		lastBidAmount  uint64
		anotherBid     uint64
		bidAmount      uint64
		editionLimit   uint64
		auctionId      uint64
		expectPass     bool
	}{
		{
			"Not existing auction id",
			types.AuctionState_Ended,
			types.AuctionPrizeType_NftOnlyTransfer,
			false,
			1500,
			0,
			1000,
			0,
			0,
			false,
		},
		{
			"claiming not existing bid",
			types.AuctionState_Ended,
			types.AuctionPrizeType_NftOnlyTransfer,
			false,
			1500,
			0,
			0,
			0,
			1,
			false,
		},
		{
			"try to claim not winner bid",
			types.AuctionState_Ended,
			types.AuctionPrizeType_NftOnlyTransfer,
			false,
			1500,
			0,
			1000,
			0,
			1,
			false,
		},
		{
			"try to claim not ended auction",
			types.AuctionState_Started,
			types.AuctionPrizeType_NftOnlyTransfer,
			false,
			1000,
			0,
			1000,
			0,
			1,
			false,
		},
		{
			"successful bid claim - with presale happened nft",
			types.AuctionState_Ended,
			types.AuctionPrizeType_NftOnlyTransfer,
			true,
			1000,
			0,
			1000,
			0,
			1,
			true,
		},
		{
			"successful bid claim - with full rights transfer",
			types.AuctionState_Ended,
			types.AuctionPrizeType_FullRightsTransfer,
			false,
			1000,
			0,
			1000,
			0,
			1,
			true,
		},
		{
			"successful bid claim - with nft only transfer",
			types.AuctionState_Ended,
			types.AuctionPrizeType_NftOnlyTransfer,
			false,
			1000,
			0,
			1000,
			0,
			1,
			true,
		},
		{
			"not successful bid claim - limited edition",
			types.AuctionState_Ended,
			types.AuctionPrizeType_LimitedEditionPrints,
			false,
			1000,
			1100,
			1000,
			1,
			1,
			false,
		},
		{
			"successful bid claim - limited edition",
			types.AuctionState_Ended,
			types.AuctionPrizeType_LimitedEditionPrints,
			false,
			1000,
			1100,
			1000,
			2,
			1,
			true,
		},
		{
			"successful bid claim - open edition",
			types.AuctionState_Ended,
			types.AuctionPrizeType_OpenEditionPrints,
			false,
			1000,
			1100,
			1000,
			0,
			1,
			true,
		},
	}

	for _, tc := range tests {
		// module account
		moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

		// set nft with ownership
		nft := nfttypes.NFT{
			CollId:     1,
			MetadataId: 1,
			Seq:        0,
			Owner:      moduleAddr.String(),
		}
		suite.app.NFTKeeper.SetNFT(suite.ctx, nft)

		// set metadata with ownership
		metadata := nfttypes.Metadata{
			Id:                   1,
			MetadataAuthority:    moduleAddr.String(),
			MintAuthority:        moduleAddr.String(),
			PrimarySaleHappened:  tc.presaleHappend,
			Name:                 "NewPUNK1",
			SellerFeeBasisPoints: 10,
			Creators: []nfttypes.Creator{
				{Address: creator.String(), Verified: true, Share: 100},
			},
			MasterEdition: &nfttypes.MasterEdition{
				Supply:    1,
				MaxSupply: 1000,
			},
		}
		suite.app.NFTKeeper.SetMetadata(suite.ctx, metadata)

		// set auction with ownership
		auction := types.Auction{
			Id:            1,
			Authority:     owner.String(),
			NftId:         nft.Id(),
			Duration:      time.Second,
			PrizeType:     tc.prizeType,
			State:         tc.state,
			LastBidAmount: tc.lastBidAmount,
			BidDenom:      "ubtsg",
			EditionLimit:  tc.editionLimit,
		}
		suite.app.MarketplaceKeeper.SetAuction(suite.ctx, auction)

		bids := suite.app.MarketplaceKeeper.GetAllBids(suite.ctx)
		for _, bid := range bids {
			suite.app.MarketplaceKeeper.DeleteBid(suite.ctx, bid)
		}

		if tc.anotherBid > 0 {
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder2.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.anotherBid,
			})
		}

		if tc.bidAmount > 0 {
			suite.app.MarketplaceKeeper.SetBidderMetadata(suite.ctx, types.BidderMetadata{
				Bidder:           bidder.String(),
				LastAuctionId:    tc.auctionId,
				LastBid:          tc.bidAmount,
				LastBidTimestamp: suite.ctx.BlockTime(),
				LastBidCancelled: false,
			})
			suite.app.MarketplaceKeeper.SetBid(suite.ctx, types.Bid{
				Bidder:    bidder.String(),
				AuctionId: tc.auctionId,
				Amount:    tc.bidAmount,
			})
		}

		oldAuctionAuthorityBalance := suite.app.BankKeeper.GetBalance(suite.ctx, owner, "ubtsg")
		oldModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, "ubtsg")
		oldCreatorBalance := suite.app.BankKeeper.GetBalance(suite.ctx, creator, "ubtsg")

		// execute SetAuctionAuthority
		msg := types.NewMsgClaimBid(bidder, tc.auctionId)
		err := suite.app.MarketplaceKeeper.ClaimBid(suite.ctx, msg)

		// check error exists on the execution
		if tc.expectPass {
			suite.Require().NoError(err)

			// check tokens are sent to auction authority from module account
			newAuctionAuthorityBalance := suite.app.BankKeeper.GetBalance(suite.ctx, owner, "ubtsg")
			newModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, "ubtsg")
			newCreatorBalance := suite.app.BankKeeper.GetBalance(suite.ctx, creator, "ubtsg")

			suite.Require().True(newAuctionAuthorityBalance.Amount.GT(oldAuctionAuthorityBalance.Amount))
			suite.Require().True(oldModuleBalance.Amount.GT(newModuleBalance.Amount))

			if tc.presaleHappend {
				// check royalties are paid if presale happened
				suite.Require().True(oldCreatorBalance.Amount.LT(newCreatorBalance.Amount))
			} else {
				// check royalties are not paid if presale not happened
				suite.Require().Equal(oldCreatorBalance, newCreatorBalance)
			}

			// check presale happened true flag is set
			newmeta, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, metadata.Id)
			suite.Require().NoError(err)
			suite.Require().True(newmeta.PrimarySaleHappened)

			// check auction Claimed value increase
			newAuction, err := suite.app.MarketplaceKeeper.GetAuctionById(suite.ctx, tc.auctionId)
			suite.Require().NoError(err)
			suite.Require().Equal(newAuction.Claimed, uint64(1))

			// check nft ownership is transfered to the bidder if nft transfer
			switch tc.prizeType {
			case types.AuctionPrizeType_FullRightsTransfer:
				// check metadata ownership is also transfered to bidder if full rights transfer auction
				suite.Require().Equal(newmeta.MetadataAuthority, bidder.String())
				fallthrough
			case types.AuctionPrizeType_NftOnlyTransfer:
				newNft, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, auction.NftId)
				suite.Require().NoError(err)
				suite.Require().Equal(newNft.Owner, bidder.String())
			case types.AuctionPrizeType_LimitedEditionPrints:
				fallthrough
			case types.AuctionPrizeType_OpenEditionPrints:
				nfts := suite.app.NFTKeeper.GetNFTsByOwner(suite.ctx, bidder)
				suite.Require().Greater(len(nfts), 0)
				suite.Require().Greater(nfts[len(nfts)-1].Seq, uint64(0))
			}

			// try claim again after successful execution
			err = suite.app.MarketplaceKeeper.ClaimBid(suite.ctx, msg)
			suite.Require().Error(err)
		} else {
			suite.Require().Error(err)
		}
	}
}
