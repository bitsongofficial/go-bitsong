package track

import (
	"fmt"
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// set the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func EndBlocker(ctx sdk.Context, k Keeper) {
	blockHeight := ctx.BlockHeight()
	blocksToPay := int64(2) // pay each 2 blocks

	if blockHeight%blocksToPay == 0 {
		// calculate reward
		// FORMULA
		// PA = Pool Amount
		// TSP = Total Streams Power
		// WSP = Weight Stream Power (PA / TSP)
		//
		// Each Track
		// U = Users
		// S = Streams
		// P = Power (Shares)
		// UP = Users Power
		// SP = Streams Power (S / U) * UP
		//
		// R = Reward
		// R = UP * WSP
		// then calculate redistribution split rate per account
		//tokens := sdk.DecCoins{{sdk.DefaultBondDenom, sdk.NewDec(1000000)}}
		// TODO: change to track owner
		//k.AllocateTokensToAccount(ctx, play.AccAddress, tokens)

		// get initial data
		playPool := k.GetPlayPool(ctx)
		playCollection := make(map[uint64]types.PlayReward)

		fmt.Println()
		fmt.Printf("-------------------")
		fmt.Println()
		fmt.Printf("Play Pool: %s", playPool)
		fmt.Println()

		// IterateAllPlays - calculate some parameters each play and add it on playCollection
		k.IterateAllPlays(ctx, func(play types.Play) bool {
			fmt.Printf("Iterate track id: %v", play.TrackId)
			fmt.Println()
			reward, ok := playCollection[play.TrackId]
			if !ok {
				reward = types.PlayReward{
					Streams:      play.Streams,
					Users:        sdk.NewInt(1),
					UsersPower:   play.Shares,
					StreamsPower: play.Streams.Quo(sdk.NewInt(1)).Mul(play.Shares.TruncateInt()),
					Reward:       sdk.NewInt64Coin("ubtsg", 0),
					Owner:        play.AccAddress,
					TrackID:      play.TrackId,
				}
				fmt.Println("Create new reward")
				fmt.Println()
			} else {
				reward.Streams = reward.Streams.Add(play.Streams)
				reward.Users = reward.Users.Add(sdk.NewInt(1))
				reward.UsersPower = reward.UsersPower.Add(play.Shares)
				reward.StreamsPower = reward.Streams.Quo(reward.Users).Mul(reward.UsersPower.TruncateInt())
				fmt.Println("Edit reward")
				fmt.Println()
			}
			playCollection[play.TrackId] = reward

			// Delete play
			k.DeletePlay(ctx, play)

			fmt.Printf("Play deleted: %v", play.TrackId)
			fmt.Println()

			return false
		})

		// Calculate total Stream Power
		totalStreamPower := sdk.ZeroInt()
		for play := range playCollection {
			totalStreamPower = totalStreamPower.Add(playCollection[play].StreamsPower)
		}

		fmt.Printf("Total Stream Power: %s", totalStreamPower)
		fmt.Println()
		fmt.Printf("Play pool rewards: %s ", playPool.Rewards)
		fmt.Println()
		fmt.Printf("Coin into module: %s", k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
		fmt.Println()
		// If totalStreamPower is > 0
		if totalStreamPower.GT(sdk.NewInt(0)) {
			// Calculate price stream power
			//priceStreamPower := playPool.Rewards.AmountOf("ubtsg").Quo(totalStreamPower.ToDec())
			priceStreamPower := playPool.Rewards.ToDec().Quo(totalStreamPower.ToDec())
			fmt.Printf("Price Stream Power: %s", priceStreamPower)
			fmt.Println()

			for play := range playCollection {
				item, ok := playCollection[play]
				if ok {
					reward := priceStreamPower.MulTruncate(item.StreamsPower.ToDec())
					coin := sdk.NewCoin("ubtsg", reward.TruncateInt())
					fmt.Printf("Reward Track %v - %s", item.TrackID, coin)
					fmt.Println()

					// Send token to track owner
					err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, item.Owner, sdk.NewCoins(coin))
					if err != nil {
						panic(err)
					}

					// Adjust play pool
					fmt.Printf("Set new play pool")
					fmt.Println()
					playPool.Rewards = playPool.Rewards.Sub(reward.TruncateInt())
					k.SetPlayPool(ctx, playPool)
					fmt.Printf("New play pool: %v", playPool.Rewards)
					fmt.Println()
				}
			}
		}

		fmt.Printf("-------------------")
		fmt.Println()
	}
}
