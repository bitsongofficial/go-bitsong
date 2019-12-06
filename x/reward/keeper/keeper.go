package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/bitsongofficial/go-bitsong/x/reward/types"

	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	paramSpace   params.Subspace
	supplyKeeper supply.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, supplyKeeper supply.Keeper) Keeper {
	// ensure distribution module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace.WithKeyTable(ParamKeyTable()),
		supplyKeeper: supplyKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetRewardPool(ctx sdk.Context) (rewardPool types.RewardPool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(RewardPoolKey)
	if b == nil {
		panic("Stored reward pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewardPool)
	return
}

func (k Keeper) SetRewardPool(ctx sdk.Context, rewardPool types.RewardPool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(rewardPool)
	store.Set(RewardPoolKey, b)
}

func (k Keeper) GetRewardModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) AddCollectedCoins(ctx sdk.Context, coins sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, "mint", types.ModuleName, coins)
}

func (k Keeper) GetRewardPoolSupply(ctx sdk.Context) sdk.Coins {
	account := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	return account.GetCoins()
}