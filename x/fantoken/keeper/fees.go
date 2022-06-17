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

// deductMintFee performs fee handling for minting token
func (k Keeper) deductMintFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// send mint fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.MintFee}, authority)
}

// deductBurnFee performs fee handling for burning token
func (k Keeper) deductBurnFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// send burn fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.BurnFee}, authority)
}

// deductTransferFee performs fee handling for transfer token
func (k Keeper) deductTransferFee(ctx sdk.Context, authority sdk.AccAddress) error {
	params := k.GetParamSet(ctx)

	// send transfer fantoken fee to community pool
	return k.distrKeeper.FundCommunityPool(ctx, sdk.Coins{params.TransferFee}, authority)
}
