package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress) (types.Bid, error) {
	bid := types.Bid{}
	auctionIdBz := sdk.Uint64ToBigEndian(auctionId)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(append(types.PrefixBid, auctionIdBz...), bidder...))
	if bz == nil {
		return bid, types.ErrBidDoesNotExists
	}

	k.cdc.MustUnmarshal(bz, &bid)
	return bid, nil
}

func (k Keeper) GetBidsByAuction(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress) []types.Bid {
	store := ctx.KVStore(k.storeKey)

	bids := []types.Bid{}
	auctionIdBz := sdk.Uint64ToBigEndian(auctionId)
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixBid, auctionIdBz...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		bid := types.Bid{}
		k.cdc.MustUnmarshal(it.Value(), &bid)

		bids = append(bids, bid)
	}
	return bids
}

func (k Keeper) GetBidsByBidder(ctx sdk.Context, bidder sdk.AccAddress) []types.Bid {
	store := ctx.KVStore(k.storeKey)

	bids := []types.Bid{}
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixBidByBidder, bidder...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		bid := types.Bid{}
		k.cdc.MustUnmarshal(it.Value(), &bid)

		bids = append(bids, bid)
	}
	return bids
}

func (k Keeper) SetBid(ctx sdk.Context, bid types.Bid) {
	// if previous bid exists, delete it
	bidder, err := sdk.AccAddressFromBech32(bid.Bidder)
	if err != nil {
		panic(err)
	}
	if bid, err := k.GetBid(ctx, bid.AuctionId, bidder); err == nil {
		k.DeleteBid(ctx, bid)
	}

	auctionIdBz := sdk.Uint64ToBigEndian(bid.AuctionId)
	bz := k.cdc.MustMarshal(&bid)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(append(types.PrefixBid, auctionIdBz...), bidder...), bz)
	store.Set(append(append(types.PrefixBidByBidder, bidder...), auctionIdBz...), bz)
}

func (k Keeper) DeleteBid(ctx sdk.Context, bid types.Bid) {
	bidder, err := sdk.AccAddressFromBech32(bid.Bidder)
	if err != nil {
		panic(err)
	}
	auctionIdBz := sdk.Uint64ToBigEndian(bid.AuctionId)
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(append(types.PrefixBid, auctionIdBz...), bidder...))
	store.Delete(append(append(types.PrefixBidByBidder, bidder...), auctionIdBz...))
}

func (k Keeper) SetBidderMetadata(ctx sdk.Context, bidderdata types.BidderMetadata) {
	bidder, err := sdk.AccAddressFromBech32(bidderdata.Bidder)
	if err != nil {
		panic(err)
	}

	bz := k.cdc.MustMarshal(&bidderdata)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixBidderMetadata, bidder...), bz)
}

func (k Keeper) GetBidderMetadata(ctx sdk.Context, bidder sdk.AccAddress) (types.BidderMetadata, error) {
	bidderdata := types.BidderMetadata{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixBidderMetadata, bidder...))
	if bz == nil {
		return bidderdata, types.ErrBidderMetadataDoesNotExists
	}

	k.cdc.MustUnmarshal(bz, &bidderdata)
	return bidderdata, nil
}

func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) error {

	// Verify auction is `Started` status
	auction, err := k.GetAuctionById(ctx, msg.AuctionId)
	if err != nil {
		return err
	}
	if auction.State != types.AuctionState_Started {
		return types.ErrAuctionNotStarted
	}

	// Verify bid is valid for the auction (check `bid_denom`, `tick_size` and `last_bid_amount`)
	if auction.BidDenom != msg.Amount.Denom {
		return types.ErrInvalidBidDenom
	}

	tickSize := auction.TickSize
	if tickSize == 0 {
		tickSize = 1
	}
	if sdk.NewInt(int64(auction.LastBidAmount+tickSize)).GT(msg.Amount.Amount) ||
		msg.Amount.Amount.LT(sdk.NewInt(int64(auction.PriceFloor))) {
		return types.ErrInvalidBidAmount
	}

	bidder, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	// check if previous bid exists for this auction by the bidder and if exists reject
	if _, err := k.GetBid(ctx, msg.AuctionId, bidder); err == nil {
		return types.ErrBidAlreadyExists
	}

	// Add new bid for the auction on the storage
	k.SetBid(ctx, types.Bid{
		AuctionId: msg.AuctionId,
		Bidder:    msg.Sender,
		Amount:    msg.Amount.Amount.Uint64(),
	})

	// Transfer amount of token to bid account
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidder, types.ModuleName, sdk.Coins{msg.Amount})
	if err != nil {
		return err
	}

	// Serialize new auction state with new bid
	auction.LastBid = ctx.BlockTime()
	auction.LastBidAmount = msg.Amount.Amount.Uint64()
	k.SetAuction(ctx, auction)

	// Update or create bidder metadata
	k.SetBidderMetadata(ctx, types.BidderMetadata{
		Bidder:           msg.Sender,
		LastAuctionId:    msg.AuctionId,
		LastBid:          msg.Amount.Amount.Uint64(),
		LastBidTimestamp: ctx.BlockTime(),
		LastBidCancelled: false,
	})

	// If the amount exceeds `instant_sale_price`, end the auction
	if msg.Amount.Amount.GTE(sdk.NewInt(int64(auction.InstantSalePrice))) {
		err := k.EndAuction(ctx, &types.MsgEndAuction{
			Sender:    auction.Authority,
			AuctionId: auction.Id,
		})
		if err != nil {
			return err
		}
	}

	// Emit event for placing bid
	ctx.EventManager().EmitTypedEvent(&types.EventPlaceBid{
		Bidder:    msg.Sender,
		AuctionId: msg.AuctionId,
	})

	return nil
}
