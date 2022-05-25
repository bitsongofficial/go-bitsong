package marketplace

import (
	"github.com/bitsongofficial/go-bitsong/x/marketplace/keeper"
	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	auctions := k.GetAuctionsToEnd(ctx)

	for _, auction := range auctions {
		authority, err := sdk.AccAddressFromBech32(auction.Authority)
		if err != nil {
			continue
		}
		cacheCtx, write := ctx.CacheContext()
		err = k.EndAuction(cacheCtx, types.NewMsgEndAuction(authority, auction.Id))
		if err == nil {
			write()
		}
	}
}
