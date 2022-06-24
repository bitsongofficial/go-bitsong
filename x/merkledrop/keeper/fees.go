package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// DeductCreationFee performs fee handling for merkledrop creation
func (k Keeper) DeductCreationFee(ctx sdk.Context, owner sdk.AccAddress) (sdk.Coin, error) {
	// send issue fantoken fee to community pool
	params := k.GetParamSet(ctx)
	return params.CreationFee, k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.CreationFee}, owner)
}
