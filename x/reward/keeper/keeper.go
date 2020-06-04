package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/reward/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	bankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper creates a content keeper
func NewKeeper(bk bank.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		bankKeeper: bk,
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetRewardPool(ctx sdk.Context, rp types.RewardPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rp)
	store.Set(types.RewardPoolKey, bz)
}

func (k Keeper) GetRewardPool(ctx sdk.Context) (rp types.RewardPool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardPoolKey)
	if bz == nil {
		panic("Stored reward pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rp)
	return rp
}

func (k Keeper) AddCollectedCoins(ctx sdk.Context, module string, coins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, module, types.ModuleName, coins)
}

func (k Keeper) GetRewardPoolSupply(ctx sdk.Context) sdk.Coins {
	return k.GetRewardPool(ctx).Amount
}
