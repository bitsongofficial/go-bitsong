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
		// Pay plays
		k.IterateAllPlays(ctx, func(play types.Play) bool {
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

			// send token to owner wallet
			playPool := k.GetFeePlayPool(ctx)

			fmt.Println()
			fmt.Printf("%s", playPool.Rewards.String())
			fmt.Println()

			err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, play.AccAddress, playPool.Rewards)
			if err != nil {
				panic(err)
			}

			// Reset play pool
			playPool.Rewards = playPool.Rewards.Sub(playPool.Rewards)
			k.SetFeePlayPool(ctx, playPool)

			// Delete play
			k.DeletePlay(ctx, play)

			return false
		})

	}
}
