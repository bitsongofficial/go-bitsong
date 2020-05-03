package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/player/types"
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

func (k Keeper) SetPlayersCount(ctx sdk.Context, count uint64) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(count)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PlayersCountKey, bz)
}

func (k Keeper) GetPlayersCount(ctx sdk.Context) (count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlayersCountKey)
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	return count
}

func (k Keeper) SetPlayer(ctx sdk.Context, player types.Player) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(player)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PlayerKey(player.ID), bz)
}

func (k Keeper) GetNode(ctx sdk.Context, id types.PlayerID) (player types.Player, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlayerKey(id))
	if bz == nil {
		return player, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &player)
	return player, true
}

func (k Keeper) AddDeposit(ctx sdk.Context, address sdk.AccAddress, coin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(coin))
}
