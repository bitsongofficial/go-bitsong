package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track"
	trackTypes "github.com/bitsongofficial/go-bitsong/x/track/types"
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
	trackKeeper  track.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, supplyKeeper supply.Keeper, trackKeeper track.Keeper) Keeper {
	// ensure distribution module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace.WithKeyTable(ParamKeyTable()),
		supplyKeeper: supplyKeeper,
		trackKeeper:  trackKeeper,
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

func (k Keeper) GetAllShares(ctx sdk.Context) trackTypes.Shares {
	return k.trackKeeper.GetAllShares(ctx)
}

func (k Keeper) GetTrack(ctx sdk.Context, trackID uint64) (track trackTypes.Track, ok bool) {
	return k.trackKeeper.GetTrack(ctx, trackID)
}

func (k Keeper) AllocateToken(ctx sdk.Context, accAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	reward, ok := k.GetReward(ctx, accAddr)

	if !ok {
		reward = types.NewReward(accAddr)
	}

	reward.TotalRewards = reward.TotalRewards.Add(amt)
	k.setReward(ctx, accAddr, reward)

	return nil
}

func (keeper Keeper) setReward(ctx sdk.Context, accAddr sdk.AccAddress, reward types.Reward) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(reward)
	store.Set(types.RewardKey(accAddr), bz)
}

func (keeper Keeper) GetReward(ctx sdk.Context, accAddr sdk.AccAddress) (reward types.Reward, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.RewardKey(accAddr))
	if bz == nil {
		return reward, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &reward)
	return reward, true
}

func (keeper Keeper) IterateAllRewards(ctx sdk.Context, cb func(reward types.Reward) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RewardsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var reward types.Reward
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reward)

		if cb(reward) {
			break
		}
	}
}

func (keeper Keeper) GetAllRewards(ctx sdk.Context) (rewards types.Rewards) {
	keeper.IterateAllRewards(ctx, func(reward types.Reward) bool {
		rewards = append(rewards, reward)
		return false
	})
	return
}
