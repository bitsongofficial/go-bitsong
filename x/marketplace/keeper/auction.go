package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
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

func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) {
	idBz := sdk.Uint64ToBigEndian(auction.Id)
	bz := k.cdc.MustMarshal(&auction)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixAuction, idBz...), bz)
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
	k.nftKeeper.TransferNFT(ctx, &nfttypes.MsgTransferNFT{
		Sender:   msg.Sender,
		Id:       msg.NftId,
		NewOwner: moduleAddr.String(),
	})

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
