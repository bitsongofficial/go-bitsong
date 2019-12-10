package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, storekey sdk.StoreKey, codespaceType sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  storekey,
		cdc:       cdc,
		codespace: codespaceType,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetDistributor(ctx sdk.Context, distributor types.Distributor) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(distributor)
	store.Set(types.DistributorKey(distributor.Address), bz)
}

func (k Keeper) GetDistributor(ctx sdk.Context, addr sdk.AccAddress) (distributor types.Distributor, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DistributorKey(addr))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &distributor)
	return distributor, true
}

func (k Keeper) CreateDistributor(ctx sdk.Context, name string, owner sdk.AccAddress) (types.Distributor, sdk.Error) {
	distributor, ok := k.GetDistributor(ctx, owner)
	if ok {
		// TODO: add duplicate distributor error
		return types.Distributor{}, nil
	}

	distributor = types.NewDistributor(name, owner)
	k.SetDistributor(ctx, distributor)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateDistributor,
			sdk.NewAttribute(types.AttributeKeyDistributorName, name),
			sdk.NewAttribute(types.AttributeKeyDistributorOwner, fmt.Sprintf("%s", owner.String())),
		),
	)

	return distributor, nil
}

func (k Keeper) IterateAllDistributors(ctx sdk.Context, cb func(distributor types.Distributor) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DistributorsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var distributor types.Distributor
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &distributor)

		if cb(distributor) {
			break
		}
	}
}

func (k Keeper) GetAllDistributors(ctx sdk.Context) (distributors types.Distributors) {
	k.IterateAllDistributors(ctx, func(distributor types.Distributor) bool {
		distributors = append(distributors, distributor)
		return false
	})
	return
}
