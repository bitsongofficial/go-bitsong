//nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// deductIssueFee performs fee handling for issuing token
func (k Keeper) deductIssueFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// check if amount is zero
	if params.IssueFee.Amount.IsZero() || params.IssueFee.Amount.IsNegative() {
		return nil
	}

	// send issue fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.IssueFee}, authority)
}

// deductMintFee performs fee handling for minting token
func (k Keeper) deductMintFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// check if amount is zero
	if params.MintFee.Amount.IsZero() || params.MintFee.Amount.IsNegative() {
		return nil
	}

	// send mint fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.MintFee}, authority)
}

// deductBurnFee performs fee handling for burning token
func (k Keeper) deductBurnFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// check if amount is zero
	if params.BurnFee.Amount.IsZero() || params.BurnFee.Amount.IsNegative() {
		return nil
	}

	// send burn fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.BurnFee}, authority)
}
