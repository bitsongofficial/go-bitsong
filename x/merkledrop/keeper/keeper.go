package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Codec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: ak,
		bankKeeper:    bk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("go-bitsong/%s", types.ModuleName))
}