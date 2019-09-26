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
		playPool := k.GetFeePlayPool(ctx)
		playCollection := make(map[uint64]types.PlayReward)

		// IterateAllPlays - calculate some parameters each play and add it on playCollection
		k.IterateAllPlays(ctx, func(play types.Play) bool {
			reward, ok := playCollection[play.TrackId]
			if !ok {
				reward = types.PlayReward{
					Streams:      play.Streams,
					Users:        sdk.NewInt(1),
					UsersPower:   play.Shares,
					StreamsPower: play.Streams.Quo(sdk.NewInt(1)).Mul(play.Shares.TruncateInt()),
					Reward:       sdk.NewInt64Coin("ubtsg", 0),
					Owner:        play.AccAddress,
				}

			} else {
				reward.Streams = reward.Streams.Add(play.Streams)
				reward.Users = reward.Users.Add(sdk.NewInt(1))
				reward.UsersPower = reward.UsersPower.Add(play.Shares)
				reward.StreamsPower = reward.Streams.Quo(reward.Users).Mul(reward.UsersPower.TruncateInt())
			}
			playCollection[play.TrackId] = reward

			// Delete play
			k.DeletePlay(ctx, play)

			return false
		})

		// Calculate total Stream Power
		totalStreamPower := sdk.ZeroInt()
		for play := range playCollection {
			totalStreamPower = totalStreamPower.Add(playCollection[play].StreamsPower)
		}

		fmt.Printf("%s", totalStreamPower)
		fmt.Println()
		fmt.Printf("Play pool reward %s ", playPool.Rewards.AmountOf("ubtsg"))
		fmt.Println()
		fmt.Printf("coin into module %s", k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
		fmt.Println()
		// If totalStreamPower is > 0
		if totalStreamPower.GT(sdk.NewInt(0)) {
			// Calculate price stream power
			priceStreamPower := playPool.Rewards.AmountOf("ubtsg").Quo(totalStreamPower).ToDec()
			fmt.Printf("%s", priceStreamPower)
			fmt.Println()

			for play := range playCollection {
				item, ok := playCollection[play]
				if ok {
					reward := sdk.NewCoin("ubtsg", priceStreamPower.MulTruncate(item.StreamsPower.ToDec()).TruncateInt())
					fmt.Printf("reward %s", reward)
					fmt.Println()

					// Send token to track owner
					err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, item.Owner, sdk.NewCoins(reward))
					if err != nil {
						panic(err)
					}

					// Adjust play pool
					playPool.Rewards = playPool.Rewards.Sub(sdk.NewCoins(reward))
					k.SetFeePlayPool(ctx, playPool)
					fmt.Printf("new play pool %v", playPool.Rewards)
				}
			}
		}
	}
}
