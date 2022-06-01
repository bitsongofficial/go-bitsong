package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

func (k Keeper) GetLastAuctionId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastAuctionId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastAuctionId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastAuctionId, idBz)
}

func (k Keeper) GetAuctionById(ctx sdk.Context, id uint64) (types.Auction, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixAuction, sdk.Uint64ToBigEndian(id)...))
	if bz == nil {
		return types.Auction{}, sdkerrors.Wrapf(types.ErrAuctionDoesNotExist, "auction: %d does not exist", id)
	}
	auction := types.Auction{}
	k.cdc.MustUnmarshal(bz, &auction)
	return auction, nil
}

func (k Keeper) GetAuctionsByAuthority(ctx sdk.Context, authority sdk.AccAddress) []types.Auction {
	store := ctx.KVStore(k.storeKey)

	auctions := []types.Auction{}
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixAuctionByAuthority, authority...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		auction, err := k.GetAuctionById(ctx, id)
		if err != nil {
			panic(err)
		}

		auctions = append(auctions, auction)
	}
	return auctions
}

func getTimeKey(timestamp time.Time) []byte {
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len(types.PrefixAuctionByEndTime)

	bz := make([]byte, prefixL+8+timeBzL)

	// copy the prefix
	copy(bz[:prefixL], types.PrefixAuctionByEndTime)

	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))

	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)
	return bz
}

func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) {
	// if previous auction exists, delete it
	if oldAuction, err := k.GetAuctionById(ctx, auction.Id); err == nil {
		k.DeleteAuction(ctx, oldAuction)
	}

	idBz := sdk.Uint64ToBigEndian(auction.Id)
	bz := k.cdc.MustMarshal(&auction)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixAuction, idBz...), bz)

	owner, err := sdk.AccAddressFromBech32(auction.Authority)
	if err != nil {
		panic(err)
	}
	store.Set(append(append(types.PrefixAuctionByAuthority, owner...), idBz...), idBz)

	if auction.IsActive() {
		store.Set(append(getTimeKey(auction.EndAuctionAt), idBz...), idBz)
	}
}

func (k Keeper) DeleteAuction(ctx sdk.Context, auction types.Auction) {
	idBz := sdk.Uint64ToBigEndian(auction.Id)
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.PrefixAuction, idBz...))

	owner, err := sdk.AccAddressFromBech32(auction.Authority)
	if err != nil {
		panic(err)
	}
	store.Delete(append(append(types.PrefixAuctionByAuthority, owner...), idBz...))

	if auction.IsActive() {
		store.Delete(append(getTimeKey(auction.EndAuctionAt), idBz...))
	}
}

func (k Keeper) GetAuctionsToEnd(ctx sdk.Context) []types.Auction {
	store := ctx.KVStore(k.storeKey)
	timeKey := getTimeKey(ctx.BlockTime())
	it := store.Iterator(types.PrefixAuctionByEndTime, storetypes.InclusiveEndBytes(timeKey))
	defer it.Close()

	auctions := []types.Auction{}
	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		auction, err := k.GetAuctionById(ctx, id)
		if err != nil {
			panic(err)
		}

		auctions = append(auctions, auction)
	}
	return auctions
}

func (k Keeper) GetAllAuctions(ctx sdk.Context) []types.Auction {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.PrefixAuction)
	defer it.Close()

	allAuctions := []types.Auction{}
	for ; it.Valid(); it.Next() {
		var auction types.Auction
		k.cdc.MustUnmarshal(it.Value(), &auction)

		allAuctions = append(allAuctions, auction)
	}

	return allAuctions
}

func (k Keeper) CreateAuction(ctx sdk.Context, msg *types.MsgCreateAuction) (uint64, error) {

	// burn fees before minting an nft
	fee := k.GetParamSet(ctx).AuctionCreationPrice
	if fee.IsPositive() {
		feeCoins := sdk.Coins{fee}
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return 0, err
		}
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, feeCoins)
		if err != nil {
			return 0, err
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, feeCoins)
		if err != nil {
			return 0, err
		}
	}

	// Ensure nft is owned by the sender
	nft, err := k.nftKeeper.GetNFTById(ctx, msg.NftId)
	if err != nil {
		return 0, err
	}
	if nft.Owner != msg.Sender {
		return 0, nfttypes.ErrNotNFTOwner
	}

	// Send nft ownership to marketplace module
	moduleAddr := k.accKeeper.GetModuleAddress(types.ModuleName)
	err = k.nftKeeper.TransferNFT(ctx, &nfttypes.MsgTransferNFT{
		Sender:   msg.Sender,
		Id:       msg.NftId,
		NewOwner: moduleAddr.String(),
	})
	if err != nil {
		return 0, err
	}

	// If auction is for transferring metadata ownership as well, metadata authority is transferred to marketplace module
	if msg.PrizeType == types.AuctionPrizeType_FullRightsTransfer {
		metadata, err := k.nftKeeper.GetMetadataById(ctx, nft.MetadataId)
		if err != nil {
			return 0, err
		}

		// Ensure nft metadata is owned by the sender if auction prize type is `FullRightsTransfer`
		if metadata.UpdateAuthority != msg.Sender {
			return 0, nfttypes.ErrNotEnoughPermission
		}
		k.nftKeeper.UpdateMetadataAuthority(ctx, &nfttypes.MsgUpdateMetadataAuthority{
			Sender:       msg.Sender,
			MetadataId:   nft.MetadataId,
			NewAuthority: moduleAddr.String(),
		})
	}

	// Create auction object from provided params
	auctionId := k.GetLastAuctionId(ctx) + 1
	k.SetLastAuctionId(ctx, auctionId)
	auction := types.Auction{
		Id:               auctionId,
		Authority:        msg.Sender,
		NftId:            msg.NftId,
		PrizeType:        msg.PrizeType,
		Duration:         msg.Duration,
		BidDenom:         msg.BidDenom,
		PriceFloor:       msg.PriceFloor,
		InstantSalePrice: msg.InstantSalePrice,
		TickSize:         msg.TickSize,
		State:            types.AuctionState_Created,
		LastBidAmount:    0,
		LastBid:          time.Time{},
		EndedAt:          time.Time{},
		EndAuctionAt:     time.Time{},
		Claimed:          false,
	}
	k.SetAuction(ctx, auction)

	// Emit event for auction creation
	ctx.EventManager().EmitTypedEvent(&types.EventCreateAuction{
		Creator:   msg.Sender,
		AuctionId: auctionId,
	})

	return auctionId, nil
}

func (k Keeper) StartAuction(ctx sdk.Context, msg *types.MsgStartAuction) error {

	// Check sender is auction authority
	auction, err := k.GetAuctionById(ctx, msg.AuctionId)
	if err != nil {
		return err
	}
	if auction.Authority != msg.Sender {
		return types.ErrNotAuctionAuthority
	}

	// Ensure auction status is `Created`
	if auction.State != types.AuctionState_Created {
		return types.ErrAuctionAlreadyStarted
	}

	// Calculate auction end time from current time and auction duration
	auction.EndAuctionAt = ctx.BlockTime().Add(auction.Duration)
	// Set the state of auction to `Started`
	auction.State = types.AuctionState_Started
	// Store updated auction into store
	k.SetAuction(ctx, auction)

	// Emit event for auction start
	ctx.EventManager().EmitTypedEvent(&types.EventStartAuction{
		AuctionId: msg.AuctionId,
	})

	return nil
}

func (k Keeper) EndAuction(ctx sdk.Context, msg *types.MsgEndAuction) error {

	// Check executor is a correct authority of the auction
	auction, err := k.GetAuctionById(ctx, msg.AuctionId)
	if err != nil {
		return err
	}
	if auction.Authority != msg.Sender {
		return types.ErrNotAuctionAuthority
	}

	// Check auction is not already ended
	if auction.State == types.AuctionState_Ended {
		return types.ErrAuctionAlreadyEnded
	}

	// Set auction end time
	auction.EndedAt = ctx.BlockTime()
	// Set auction status as ended
	auction.State = types.AuctionState_Ended
	// Set updated auction on the storage
	k.SetAuction(ctx, auction)

	// If winning bid does not exists, send nft and metadata ownership back to auction authority
	if auction.LastBidAmount == 0 {
		moduleAddr := k.accKeeper.GetModuleAddress(types.ModuleName)
		k.nftKeeper.TransferNFT(ctx, &nfttypes.MsgTransferNFT{
			Sender:   moduleAddr.String(),
			Id:       auction.NftId,
			NewOwner: auction.Authority,
		})
		if auction.PrizeType == types.AuctionPrizeType_FullRightsTransfer {
			nft, err := k.nftKeeper.GetNFTById(ctx, auction.NftId)
			if err != nil {
				return err
			}

			k.nftKeeper.UpdateMetadataAuthority(ctx, &nfttypes.MsgUpdateMetadataAuthority{
				Sender:       moduleAddr.String(),
				MetadataId:   nft.MetadataId,
				NewAuthority: auction.Authority,
			})
		}
	}

	// Emit event for auction end
	ctx.EventManager().EmitTypedEvent(&types.EventEndAuction{
		AuctionId: auction.Id,
	})

	return nil
}

func (k Keeper) SetAuctionAuthority(ctx sdk.Context, msg *types.MsgSetAuctionAuthority) error {

	// Check sender is auction authority
	auction, err := k.GetAuctionById(ctx, msg.AuctionId)
	if err != nil {
		return err
	}
	if auction.Authority != msg.Sender {
		return types.ErrNotAuctionAuthority
	}

	// Ensure new authority is an accurate address
	_, err = sdk.AccAddressFromBech32(msg.NewAuthority)
	if err != nil {
		return err
	}

	// Update auction authority with new authority
	auction.Authority = msg.NewAuthority
	k.SetAuction(ctx, auction)

	// Emit event for authority update
	ctx.EventManager().EmitTypedEvent(&types.EventSetAuctionAuthority{
		AuctionId: msg.AuctionId,
		Authority: auction.Authority,
	})

	return nil
}
