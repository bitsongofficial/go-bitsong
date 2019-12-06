package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// Calculate BitSong Reward Pool
	rewardFraction, _ := sdk.NewDecFromStr("0.30") // TODO: get from parameters
	rewardCoins, _ := sdk.NewDecCoins(mintedCoins).MulDecTruncate(rewardFraction).TruncateDecimal() // truncate decimals
	remainingCoins := mintedCoins.Sub(rewardCoins) // subtract artistPool from mintedCoins

	// TODO:
	// Add artistCoins to the rewardPool
	// k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)

	err = k.AddCollectedFees(ctx, remainingCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeMint,
			sdk.NewAttribute(AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}