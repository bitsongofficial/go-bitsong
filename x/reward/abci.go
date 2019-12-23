package reward

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, keeper Keeper) {
	if ctx.BlockHeight()%2 != 1 {
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

		fmt.Println()
		fmt.Println()
		fmt.Printf("Track ID: %d", t.TrackID)
		fmt.Println()
		fmt.Printf("Reward Pool Supply: %s", rewardPoolSupply.String())
		fmt.Println()
		fmt.Printf("Reward Pool: %s", rewardPool.Amount.String())
		fmt.Println()
		fmt.Printf("Reward Portion: %s", rewardPortion.String())
		fmt.Println()
		fmt.Printf("Reward Track: %s", reward.String())
		fmt.Println()
		fmt.Printf("TotalShare Track: %s", t.TotalShare.TruncateInt().String())
		fmt.Println()
		fmt.Printf("TotalShare: %s", totalShares.TruncateInt().String())
		fmt.Println()
		fmt.Println()

		// subtract reward from rewardPool storage
		rewardPool.Amount = rewardPool.Amount.Sub(sdk.NewDecCoins(reward))
		keeper.SetRewardPool(ctx, rewardPool)
	}

	// delete all plays and shares
	keeper.DeleteAllPlays(ctx)
	keeper.DeleteAllShares(ctx)
}
