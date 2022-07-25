package candymachine

import (
	"github.com/bitsongofficial/go-bitsong/x/candymachine/keeper"
	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	machines := k.GetCandyMachinesToEndByTime(ctx)

	for _, machine := range machines {
		authority, err := sdk.AccAddressFromBech32(machine.Authority)
		if err != nil {
			continue
		}
		cacheCtx, write := ctx.CacheContext()
		err = k.CloseCandyMachine(cacheCtx, types.NewMsgCloseCandyMachine(authority))
		if err == nil {
			write()
		}
	}
}
