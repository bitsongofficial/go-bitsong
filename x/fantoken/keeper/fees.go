//nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeductIssueFanTokenFee performs fee handling for issuing token
func (k Keeper) DeductIssueFanTokenFee(ctx sdk.Context, owner sdk.AccAddress) error {
	// send issue fantoken fee to community pool
	params := k.GetParamSet(ctx)
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.IssueFee}, owner)
}
