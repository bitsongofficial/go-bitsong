package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

func (k Keeper) GetCandyMachineByCollId(ctx sdk.Context, collId uint64) (types.CandyMachine, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixCandyMachine, sdk.Uint64ToBigEndian(collId)...))
	if bz == nil {
		return types.CandyMachine{}, sdkerrors.Wrapf(types.ErrCandyMachineDoesNotExist, "candymachine: %d does not exist", collId)
	}
	candymachine := types.CandyMachine{}
	k.cdc.MustUnmarshal(bz, &candymachine)
	return candymachine, nil
}

func getTimeKey(timestamp uint64) []byte {
	time := time.Unix(int64(timestamp), 0)
	timeBz := sdk.FormatTimeBytes(time)
	timeBzL := len(timeBz)
	prefixL := len(types.PrefixCandyMachineByEndTime)

	bz := make([]byte, prefixL+8+timeBzL)

	// copy the prefix
	copy(bz[:prefixL], types.PrefixCandyMachineByEndTime)

	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))

	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)
	return bz
}

func (k Keeper) SetCandyMachine(ctx sdk.Context, machine types.CandyMachine) {
	// if previous candymachine exists, delete it
	if oldMachine, err := k.GetCandyMachineByCollId(ctx, machine.CollId); err == nil {
		k.DeleteCandyMachine(ctx, oldMachine)
	}

	idBz := sdk.Uint64ToBigEndian(machine.CollId)
	bz := k.cdc.MustMarshal(&machine)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixCandyMachine, idBz...), bz)

	if machine.EndTimestamp != 0 {
		store.Set(append(getTimeKey(machine.EndTimestamp), idBz...), idBz)
	}
}

func (k Keeper) DeleteCandyMachine(ctx sdk.Context, machine types.CandyMachine) {
	idBz := sdk.Uint64ToBigEndian(machine.CollId)
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.PrefixCandyMachine, idBz...))

	if machine.EndTimestamp != 0 {
		store.Delete(append(getTimeKey(machine.EndTimestamp), idBz...))
	}
}

func (k Keeper) GetCandyMachinesToEndByTime(ctx sdk.Context) []types.CandyMachine {
	store := ctx.KVStore(k.storeKey)
	timeKey := getTimeKey(uint64(ctx.BlockTime().Unix()))
	it := store.Iterator(types.PrefixCandyMachineByEndTime, storetypes.InclusiveEndBytes(timeKey))
	defer it.Close()

	machines := []types.CandyMachine{}
	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		machine, err := k.GetCandyMachineByCollId(ctx, id)
		if err != nil {
			panic(err)
		}

		machines = append(machines, machine)
	}
	return machines
}

func (k Keeper) GetAllCandyMachines(ctx sdk.Context) []types.CandyMachine {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.PrefixCandyMachine)
	defer it.Close()

	allMachines := []types.CandyMachine{}
	for ; it.Valid(); it.Next() {
		var machine types.CandyMachine
		k.cdc.MustUnmarshal(it.Value(), &machine)

		allMachines = append(allMachines, machine)
	}

	return allMachines
}

func (k Keeper) CreateCandyMachine(ctx sdk.Context, msg *types.MsgCreateCandyMachine) error {
	// burn fees before creating candy machine
	fee := k.GetParamSet(ctx).CandymachineCreationPrice
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
	collection, err := k.nftKeeper.GetCollectionById(ctx, msg.Machine.CollId)
	if err != nil {
		return err
	}

	if collection.UpdateAuthority != msg.Sender {
		return types.ErrNotCollectionAuthority
	}

	// max mint value limitation check
	params := k.GetParamSet(ctx)
	if msg.Machine.MaxMint > params.CandymachineMaxMint {
		return types.ErrCannotExceedMaxMintParameter
	}

	moduleAddr := k.accKeeper.GetModuleAddress(types.ModuleName)
	collection.UpdateAuthority = moduleAddr.String()
	k.nftKeeper.SetCollection(ctx, collection)

	mintableMetadataIds := []uint64{}
	lastMetadataId := k.nftKeeper.GetLastMetadataId(ctx)
	for i := uint64(0); i < msg.Machine.MaxMint; i++ {
		mintableMetadataIds = append(mintableMetadataIds, lastMetadataId+1+i)
	}
	k.nftKeeper.SetLastMetadataId(ctx, lastMetadataId+msg.Machine.MaxMint)
	mintableMetadataIds = RandomList(ctx, mintableMetadataIds)
	k.SetMintableMetadataIds(ctx, msg.Machine.CollId, mintableMetadataIds)

	machine := msg.Machine
	k.SetCandyMachine(ctx, machine)

	// Emit event for auction creation
	ctx.EventManager().EmitTypedEvent(&types.EventCreateCandyMachine{
		Creator:      msg.Sender,
		CollectionId: msg.Machine.CollId,
	})

	return nil
}

func (k Keeper) UpdateCandyMachine(ctx sdk.Context, msg *types.MsgUpdateCandyMachine) error {
	// Ensure candy machine is owned by the sender
	machine, err := k.GetCandyMachineByCollId(ctx, msg.Machine.CollId)
	if err != nil {
		return err
	}

	if machine.Authority != msg.Sender {
		return types.ErrNotCandyMachineAuthority
	}

	// TODO: make changes on machine from previous status

	params := k.GetParamSet(ctx)
	if msg.Machine.MaxMint > params.CandymachineMaxMint {
		return types.ErrCannotExceedMaxMintParameter
	}

	// if max value is increased allocate more metadata ids
	if machine.MaxMint < msg.Machine.MaxMint {
		mintableMetadataIds := k.GetMintableMetadataIds(ctx, msg.Machine.CollId)
		lastMetadataId := k.nftKeeper.GetLastMetadataId(ctx)
		for i := uint64(0); i < msg.Machine.MaxMint-machine.MaxMint; i++ {
			mintableMetadataIds = append(mintableMetadataIds, lastMetadataId+1+i)
		}
		k.nftKeeper.SetLastMetadataId(ctx, lastMetadataId+msg.Machine.MaxMint)
		mintableMetadataIds = RandomList(ctx, mintableMetadataIds)
		k.SetMintableMetadataIds(ctx, msg.Machine.CollId, mintableMetadataIds)
	}

	k.SetCandyMachine(ctx, msg.Machine)

	// Emit event for auction creation
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateCandyMachine{
		Creator:      msg.Sender,
		CollectionId: msg.Machine.CollId,
	})

	return nil
}

func (k Keeper) CloseCandyMachine(ctx sdk.Context, msg *types.MsgCloseCandyMachine) error {
	// Ensure candy machine is owned by the sender
	machine, err := k.GetCandyMachineByCollId(ctx, msg.CollId)
	if err != nil {
		return err
	}

	if machine.Authority != msg.Sender {
		return types.ErrNotCandyMachineAuthority
	}

	// delete candy machine
	k.DeleteCandyMachine(ctx, machine)
	// remove mintable metadata ids
	k.DeleteMintableMetadataIds(ctx, machine.CollId)

	// transfer ownership of collection to the sender
	collection, err := k.nftKeeper.GetCollectionById(ctx, msg.CollId)
	if err != nil {
		return err
	}

	collection.UpdateAuthority = msg.Sender
	k.nftKeeper.SetCollection(ctx, collection)

	// Emit event for candy machine close
	ctx.EventManager().EmitTypedEvent(&types.EventCloseCandyMachine{
		Creator:      msg.Sender,
		CollectionId: msg.CollId,
	})

	return nil
}

func (k Keeper) PayCandyMachineFee(ctx sdk.Context, sender sdk.AccAddress, machine types.CandyMachine) error {
	if machine.Price > 0 {
		feeCoins := sdk.Coins{sdk.NewInt64Coin(machine.Denom, int64(machine.Price))}
		authority, err := sdk.AccAddressFromBech32(machine.Authority)
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

func (k Keeper) MintNFT(ctx sdk.Context, msg *types.MsgMintNFT) error {
	// Ensure candy machine is owned by the sender
	machine, err := k.GetCandyMachineByCollId(ctx, msg.CollectionId)
	if err != nil {
		return err
	}

	if machine.GoLiveDate > uint64(ctx.BlockTime().Unix()) {
		return types.ErrCandyMachineNotLiveTime
	}

	// make payment
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	err = k.PayCandyMachineFee(ctx, sender, machine)
	if err != nil {
		return err
	}

	shuffledId := k.TakeOutRandomMintableMetadataId(ctx, machine.CollId, machine.MaxMint-machine.Minted)

	// metadata should be dynamically created on candymachine with shuffled id
	k.nftKeeper.SetMetadata(ctx, nfttypes.Metadata{
		Id:                   shuffledId,
		Name:                 msg.Name,
		Uri:                  fmt.Sprintf("%s/%d", machine.MetadataBaseUrl, machine.Minted+1),
		SellerFeeBasisPoints: machine.SellerFeeBasisPoints,
		PrimarySaleHappened:  true,
		IsMutable:            machine.Mutable,
		Creators:             machine.Creators,
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
		MetadataId: shuffledId,
		Seq:        0,
	}
	k.nftKeeper.SetNFT(ctx, nft)

	machine.Minted++

	if machine.Minted >= machine.MaxMint {
		authority, err := sdk.AccAddressFromBech32(machine.Authority)
		if err != nil {
			return err
		}
		err = k.CloseCandyMachine(ctx, types.NewMsgCloseCandyMachine(authority, machine.CollId))
		if err != nil {
			return err
		}
	} else {
		k.SetCandyMachine(ctx, machine)
	}

	// Emit event for candy machine close
	ctx.EventManager().EmitTypedEvent(&types.EventMintNFT{
		CollectionId: msg.CollectionId,
		NftId:        nft.Id(),
	})

	return nil
}

func (k Keeper) AllMintableMetadataIds(ctx sdk.Context) []types.MintableMetadataIds {
	candymachines := k.GetAllCandyMachines(ctx)

	mintableMetadataIds := []types.MintableMetadataIds{}
	for _, machine := range candymachines {
		ids := k.GetMintableMetadataIds(ctx, machine.CollId)
		mintableMetadataIds = append(mintableMetadataIds, types.MintableMetadataIds{
			CollectionId:        machine.CollId,
			MintableMetadataIds: ids,
		})
	}
	return mintableMetadataIds
}
