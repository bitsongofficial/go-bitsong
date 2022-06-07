//nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// deductIssueFee performs fee handling for issuing token
func (k Keeper) deductIssueFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// send issue fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.IssueFee}, authority)
}
