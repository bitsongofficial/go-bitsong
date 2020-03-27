package reward

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, keeper Keeper) {
	blockPeriodPayout := int64(60) // Payout each x block

	if ctx.BlockHeight()%blockPeriodPayout != 0 {
		return
	}

	logger := keeper.Logger(ctx)

	logger.Info(
		fmt.Sprintf(
			"reward endblocker",
		),
	)

	// fetch stored shares
	shares := keeper.GetAllShares(ctx)
	totalShares := sdk.NewDec(0)

	for _, t := range shares {
		totalShares = totalShares.Add(t.TotalShare)
	}

	if !totalShares.IsPositive() {
		return
	}

	// fetch rewardPool
	rewardPoolSupply := keeper.GetRewardPoolSupply(ctx)
	rewardPool := keeper.GetRewardPool(ctx)

	// calculate reward portion
	// todo: change ubtsg
	rewardPortion := rewardPool.Amount.QuoDec(totalShares)

	for _, t := range shares {
		// calculate reward
		reward, _ := rewardPortion.MulDec(t.TotalShare).TruncateDecimal()

		// get track
		track, ok := keeper.GetTrack(ctx, t.TrackID)
		if !ok {
			panic("owner not found")
		}

		// allocate reward
		keeper.AllocateToken(ctx, track, reward)

		fmt.Printf(`

Track ID: %d
Reward Pool Supply: %s
Reward Pool: %s
Reward Portion: %s
Reward Track: %s
TotalShare Track: %s
TotalShare: %s

`, t.TrackID, rewardPoolSupply, rewardPool.Amount, rewardPortion,
			reward, t.TotalShare.TruncateInt(), totalShares.TruncateInt())

		// subtract reward from rewardPool storage
		rewardPool.Amount = rewardPool.Amount.Sub(sdk.NewDecCoinsFromCoins(reward...))
		keeper.SetRewardPool(ctx, rewardPool)
	}

	// delete all plays and shares
	keeper.DeleteAllPlays(ctx)
	keeper.DeleteAllShares(ctx)
}
