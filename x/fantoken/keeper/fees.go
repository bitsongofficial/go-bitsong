//nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// DeductIssueTokenFee performs fee handling for issuing token
func (k Keeper) DeductIssueFanTokenFee(ctx sdk.Context, owner sdk.AccAddress, issueFee sdk.Coin, symbol string) error {
	burnCoins := sdk.NewCoins(issueFee)

	// send all fees to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, owner, types.ModuleName, sdk.NewCoins(issueFee),
	); err != nil {
		return err
	}

	// burn burnedCoin
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
}
