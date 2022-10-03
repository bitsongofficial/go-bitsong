package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

func (k Keeper) GetLaunchPadByCollId(ctx sdk.Context, collId uint64) (types.LaunchPad, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LaunchPadKey(collId))
	if bz == nil {
		return types.LaunchPad{}, sdkerrors.Wrapf(types.ErrLaunchPadDoesNotExist, "launchpad: %d does not exist", collId)
	}
	launchpad := types.LaunchPad{}
	k.cdc.MustUnmarshal(bz, &launchpad)
	return launchpad, nil
}

func (k Keeper) SetLaunchPad(ctx sdk.Context, pad types.LaunchPad) {
	// if previous launchpad exists, delete it
	if oldPad, err := k.GetLaunchPadByCollId(ctx, pad.CollId); err == nil {
		k.DeleteLaunchPad(ctx, oldPad)
	}

	idBz := sdk.Uint64ToBigEndian(pad.CollId)
	bz := k.cdc.MustMarshal(&pad)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LaunchPadKey(pad.CollId), bz)

	if pad.EndTimestamp != 0 {
		store.Set(types.LaunchPadByEndTimeKey(pad.EndTimestamp, pad.CollId), idBz)
	}
}

func (k Keeper) DeleteLaunchPad(ctx sdk.Context, pad types.LaunchPad) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.LaunchPadKey(pad.CollId))

	if pad.EndTimestamp != 0 {
		store.Delete(types.LaunchPadByEndTimeKey(pad.EndTimestamp, pad.CollId))
	}
}

func (k Keeper) GetLaunchPadsToEndByTime(ctx sdk.Context) []types.LaunchPad {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.GetTimeKey(uint64(ctx.BlockTime().Unix()))
	it := store.Iterator(types.PrefixLaunchPadByEndTime, storetypes.InclusiveEndBytes(timeKey))
	defer it.Close()

	pads := []types.LaunchPad{}
	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		pad, err := k.GetLaunchPadByCollId(ctx, id)
		if err != nil {
			panic(err)
		}

		pads = append(pads, pad)
	}
	return pads
}

func (k Keeper) GetAllLaunchPads(ctx sdk.Context) []types.LaunchPad {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.PrefixLaunchPad)
	defer it.Close()

	allPads := []types.LaunchPad{}
	for ; it.Valid(); it.Next() {
		var pad types.LaunchPad
		k.cdc.MustUnmarshal(it.Value(), &pad)

		allPads = append(allPads, pad)
	}

	return allPads
}

func (k Keeper) CreateLaunchPad(ctx sdk.Context, msg *types.MsgCreateLaunchPad) error {
	// burn fees before creating launchpad
	fee := k.GetParamSet(ctx).LaunchpadCreationPrice
	if fee.IsPositive() {
		feeCoins := sdk.Coins{fee}
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, feeCoins)
		if err != nil {
			return err
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, feeCoins)
		if err != nil {
			return err
		}
	}

	// Ensure nft is owned by the sender
	collection, err := k.nftKeeper.GetCollectionById(ctx, msg.Pad.CollId)
	if err != nil {
		return err
	}

	if collection.UpdateAuthority != msg.Sender {
		return types.ErrNotCollectionAuthority
	}

	// max mint value limitation check
	params := k.GetParamSet(ctx)
	if msg.Pad.MaxMint > params.LaunchpadMaxMint {
		return types.ErrCannotExceedMaxMintParameter
	}

	moduleAddr := k.accKeeper.GetModuleAddress(types.ModuleName)
	collection.UpdateAuthority = moduleAddr.String()
	k.nftKeeper.SetCollection(ctx, collection)

	mintableMetadataIds := []uint64{}
	lastMetadataId := k.nftKeeper.GetLastMetadataId(ctx, msg.Pad.CollId)
	for i := uint64(0); i < msg.Pad.MaxMint; i++ {
		mintableMetadataIds = append(mintableMetadataIds, lastMetadataId+1+i)
	}
	k.nftKeeper.SetLastMetadataId(ctx, msg.Pad.CollId, lastMetadataId+msg.Pad.MaxMint)
	if msg.Pad.Shuffle {
		mintableMetadataIds = RandomList(ctx, mintableMetadataIds)
	}
	k.SetMintableMetadataIds(ctx, msg.Pad.CollId, mintableMetadataIds)

	pad := msg.Pad
	k.SetLaunchPad(ctx, pad)

	// Emit event for auction creation
	ctx.EventManager().EmitTypedEvent(&types.EventCreateLaunchPad{
		Creator:      msg.Sender,
		CollectionId: msg.Pad.CollId,
	})

	return nil
}

func (k Keeper) UpdateLaunchPad(ctx sdk.Context, msg *types.MsgUpdateLaunchPad) error {
	// Ensure launchpad is owned by the sender
	pad, err := k.GetLaunchPadByCollId(ctx, msg.Pad.CollId)
	if err != nil {
		return err
	}

	if pad.Authority != msg.Sender {
		return types.ErrNotLaunchPadAuthority
	}

	// TODO: make changes on pad from previous status

	params := k.GetParamSet(ctx)
	if msg.Pad.MaxMint > params.LaunchpadMaxMint {
		return types.ErrCannotExceedMaxMintParameter
	}

	// if max value is increased allocate more metadata ids
	if pad.MaxMint < msg.Pad.MaxMint {
		mintableMetadataIds := k.GetMintableMetadataIds(ctx, msg.Pad.CollId)
		lastMetadataId := k.nftKeeper.GetLastMetadataId(ctx, msg.Pad.CollId)
		for i := uint64(0); i < msg.Pad.MaxMint-pad.MaxMint; i++ {
			mintableMetadataIds = append(mintableMetadataIds, lastMetadataId+1+i)
		}
		k.nftKeeper.SetLastMetadataId(ctx, msg.Pad.CollId, lastMetadataId+msg.Pad.MaxMint)
		if msg.Pad.Shuffle {
			mintableMetadataIds = RandomList(ctx, mintableMetadataIds)
		}
		k.SetMintableMetadataIds(ctx, msg.Pad.CollId, mintableMetadataIds)
	}

	k.SetLaunchPad(ctx, msg.Pad)

	// Emit event for auction creation
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateLaunchPad{
		Creator:      msg.Sender,
		CollectionId: msg.Pad.CollId,
	})

	return nil
}

func (k Keeper) CloseLaunchPad(ctx sdk.Context, msg *types.MsgCloseLaunchPad) error {
	// Ensure launchpad is owned by the sender
	pad, err := k.GetLaunchPadByCollId(ctx, msg.CollId)
	if err != nil {
		return err
	}

	if pad.Authority != msg.Sender {
		return types.ErrNotLaunchPadAuthority
	}

	// delete launchpad
	k.DeleteLaunchPad(ctx, pad)
	// remove mintable metadata ids
	k.DeleteMintableMetadataIds(ctx, pad.CollId)

	// transfer ownership of collection to the sender
	collection, err := k.nftKeeper.GetCollectionById(ctx, msg.CollId)
	if err != nil {
		return err
	}

	collection.UpdateAuthority = msg.Sender
	k.nftKeeper.SetCollection(ctx, collection)

	// Emit event for launchpad close
	ctx.EventManager().EmitTypedEvent(&types.EventCloseLaunchPad{
		Creator:      msg.Sender,
		CollectionId: msg.CollId,
	})

	return nil
}

func (k Keeper) PayLaunchPadFee(ctx sdk.Context, sender sdk.AccAddress, pad types.LaunchPad) error {
	if pad.Price > 0 {
		feeCoins := sdk.Coins{sdk.NewInt64Coin(pad.Denom, int64(pad.Price))}
		authority, err := sdk.AccAddressFromBech32(pad.Authority)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoins(ctx, sender, authority, feeCoins)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) MintNFT(ctx sdk.Context, msg *types.MsgMintNFT) (string, error) {
	// Ensure launchpad is owned by the sender
	pad, err := k.GetLaunchPadByCollId(ctx, msg.CollectionId)
	if err != nil {
		return "", err
	}

	if pad.GoLiveDate > uint64(ctx.BlockTime().Unix()) {
		return "", types.ErrLaunchPadNotLiveTime
	}

	// make payment
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return "", err
	}
	err = k.PayLaunchPadFee(ctx, sender, pad)
	if err != nil {
		return "", err
	}

	metadataId := uint64(0)
	if !pad.Shuffle {
		metadataId = k.TakeOutFirstMintableMetadataId(ctx, pad.CollId)
	} else {
		metadataId = k.TakeOutRandomMintableMetadataId(ctx, pad.CollId, pad.MaxMint-pad.Minted)
	}

	// metadata should be dynamically created on launchpad with selected metadata id
	k.nftKeeper.SetMetadata(ctx, nfttypes.Metadata{
		CollId:               msg.CollectionId,
		Id:                   metadataId,
		Name:                 msg.Name,
		Uri:                  fmt.Sprintf("%s/%d", pad.MetadataBaseUrl, pad.Minted+1),
		SellerFeeBasisPoints: pad.SellerFeeBasisPoints,
		PrimarySaleHappened:  true,
		IsMutable:            pad.Mutable,
		Creators:             pad.Creators,
		MetadataAuthority:    msg.Sender,
		MintAuthority:        msg.Sender,
		MasterEdition: &nfttypes.MasterEdition{
			Supply:    1,
			MaxSupply: 1,
		},
	})

	// create nft
	nft := nfttypes.NFT{
		Owner:      msg.Sender,
		CollId:     msg.CollectionId,
		MetadataId: metadataId,
		Seq:        0,
	}
	k.nftKeeper.SetNFT(ctx, nft)

	pad.Minted++

	if pad.Minted >= pad.MaxMint {
		authority, err := sdk.AccAddressFromBech32(pad.Authority)
		if err != nil {
			return "", err
		}
		err = k.CloseLaunchPad(ctx, types.NewMsgCloseLaunchPad(authority, pad.CollId))
		if err != nil {
			return "", err
		}
	} else {
		k.SetLaunchPad(ctx, pad)
	}

	// Emit event for launchpad close
	ctx.EventManager().EmitTypedEvent(&types.EventMintNFT{
		CollectionId: msg.CollectionId,
		NftId:        nft.Id(),
	})

	return nft.Id(), nil
}

func (k Keeper) MintNFTs(ctx sdk.Context, msg *types.MsgMintNFTs) ([]string, error) {
	collection, err := k.nftKeeper.GetCollectionById(ctx, msg.CollectionId)
	if err != nil {
		return []string{}, err
	}

	mintabeMetadataIds := k.GetMintableMetadataIds(ctx, msg.CollectionId)
	if len(mintabeMetadataIds) < int(msg.Number) {
		return []string{}, types.ErrInsufficientMintableNftsRemaining
	}

	nftIds := []string{}
	for i := uint64(1); i <= msg.Number; i++ {
		nftId, err := k.MintNFT(ctx, &types.MsgMintNFT{
			Sender:       msg.Sender,
			CollectionId: msg.CollectionId,
			Name:         fmt.Sprintf("%s #%d", collection.Name, i),
		})
		if err != nil {
			return []string{}, err
		}

		nftIds = append(nftIds, nftId)
	}
	return nftIds, nil
}

func (k Keeper) AllMintableMetadataIds(ctx sdk.Context) []types.MintableMetadataIds {
	launchpads := k.GetAllLaunchPads(ctx)

	mintableMetadataIds := []types.MintableMetadataIds{}
	for _, pad := range launchpads {
		ids := k.GetMintableMetadataIds(ctx, pad.CollId)
		mintableMetadataIds = append(mintableMetadataIds, types.MintableMetadataIds{
			CollectionId:        pad.CollId,
			MintableMetadataIds: ids,
		})
	}
	return mintableMetadataIds
}
