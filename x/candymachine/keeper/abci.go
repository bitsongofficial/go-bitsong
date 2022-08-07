package keeper

import (
	"fmt"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	machines := k.GetCandyMachinesToEndByTime(ctx)

	for _, machine := range machines {
		authority, err := sdk.AccAddressFromBech32(machine.Authority)
		if err != nil {
			continue
		}
		cacheCtx, write := ctx.CacheContext()
		err = k.CloseCandyMachine(cacheCtx, types.NewMsgCloseCandyMachine(authority, machine.CollId))
		if err == nil {
			write()
		} else {
			fmt.Println(err)
		}
	}
}
