package keeper

import (
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"cosmossdk.io/log"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

// Keeper of the cadance store
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	wasmKeeper     wasmkeeper.Keeper
	contractKeeper *wasmkeeper.PermissionedKeeper

	authority string
}

func NewKeeper(
	key storetypes.StoreKey,
	cdc codec.BinaryCodec,
	wasmKeeper wasmkeeper.Keeper,
	contractKeeper *wasmkeeper.PermissionedKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       key,
		wasmKeeper:     wasmKeeper,
		contractKeeper: contractKeeper,
		authority:      authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("go-bitsong/%s", types.ModuleName))
}

// GetAuthority returns the x/cadance module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// SetParams sets the x/cadance module parameters.
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&p)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetParams returns the current x/cadance module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

// GetContractKeeper returns the x/wasm module's contract keeper.
func (k Keeper) GetContractKeeper() *wasmkeeper.PermissionedKeeper {
	return k.contractKeeper
}

// GetCdc returns the x/cadance module's codec.
func (k Keeper) GetCdc() codec.BinaryCodec {
	return k.cdc
}

// GetStore returns the x/cadance module's store key.
func (k Keeper) GetStore() storetypes.StoreKey {
	return k.storeKey
}
