//nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeductIssueFanTokenFee performs fee handling for issuing token
func (k Keeper) DeductIssueFanTokenFee(ctx sdk.Context, owner sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// send issue fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.IssueFee}, owner)
}
