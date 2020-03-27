package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, cmintKeeper mint.Keeper, k Keeper) {
	// fetch stored minter & params
	minter := cmintKeeper.GetMinter(ctx)
	params := cmintKeeper.GetParams(ctx)

	// recalculate inflation rate
	totalStakingSupply := cmintKeeper.StakingTokenSupply(ctx)
	bondedRatio := cmintKeeper.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	cmintKeeper.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := cmintKeeper.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// Calculate BitSong Reward Pool
	rewardFraction, _ := sdk.NewDecFromStr("0.03")                                                          // TODO: (3%) get from parameters
	rewardCoins, _ := sdk.NewDecCoinsFromCoins(mintedCoin).MulDecTruncate(rewardFraction).TruncateDecimal() // truncate decimals

	// TODO:
	// Add rewardCoins to the rewardPool
	err = k.AddToRewardPool(ctx, rewardCoins)
	if err != nil {
		panic(err)
	}

	fmt.Printf(`

Reward Pool: %s

`, k.GetRewardPoolSupply(ctx))

	remainingCoins := mintedCoins.Sub(rewardCoins) // subtract artistPool from mintedCoins

	// send the minted coins to the fee collector account
	err = cmintKeeper.AddCollectedFees(ctx, remainingCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
