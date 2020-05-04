package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/player/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	bankKeeper    bank.Keeper
	stakingKeeper staking.Keeper
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
}

// NewKeeper creates a content keeper
func NewKeeper(bk bank.Keeper, sk staking.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		bankKeeper:    bk,
		stakingKeeper: sk,
		storeKey:      key,
		cdc:           cdc,
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
	store.Set(types.PlayerKey(player.PlayerAddr), bz)
}

func (k Keeper) GetPlayer(ctx sdk.Context, playerAddr sdk.AccAddress) (player types.Player, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlayerKey(playerAddr))
	if bz == nil {
		return player, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &player)
	return player, true
}

func (k Keeper) AddDeposit(ctx sdk.Context, address sdk.AccAddress, coin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.Coins{coin})
}

func (k Keeper) RegisterPlayer(ctx sdk.Context, moniker string, playerAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	val, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not a valid validator", valAddr)
	}

	if !val.IsBonded() {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to register a new player", valAddr)
	}

	_, found = k.GetPlayer(ctx, playerAddr)
	if found {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "player is already registered")
	}

	player := types.Player{
		Moniker:    moniker,
		PlayerAddr: playerAddr,
		Validator:  valAddr,
	}

	k.SetPlayer(ctx, player)

	return nil
}
