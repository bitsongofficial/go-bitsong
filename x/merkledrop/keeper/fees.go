package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// DeductCreationFee performs fee handling for merkledrop creation
func (k Keeper) DeductCreationFee(ctx sdk.Context, owner sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// check if amount is zero
	if params.CreationFee.Amount.IsZero() || params.CreationFee.Amount.IsNegative() {
		return nil
	}

	// send issue fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.CreationFee}, owner)
}
