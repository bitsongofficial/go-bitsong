package keeper

import (
	"fmt"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	pads := k.GetLaunchPadsToEndByTime(ctx)

	for _, pad := range pads {
		authority, err := sdk.AccAddressFromBech32(pad.Authority)
		if err != nil {
			continue
		}
		cacheCtx, write := ctx.CacheContext()
		err = k.CloseLaunchPad(cacheCtx, types.NewMsgCloseLaunchPad(authority, pad.CollId))
		if err == nil {
			write()
		} else {
			fmt.Println(err)
		}
	}
}
