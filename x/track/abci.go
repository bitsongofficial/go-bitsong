package track

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// set the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func EndBlocker(ctx sdk.Context, k Keeper) {
	blockHeight := ctx.BlockHeight()
	blocksToPay := int64(2)

	if blockHeight%blocksToPay == 0 {
		// Pay plays

	}
}
