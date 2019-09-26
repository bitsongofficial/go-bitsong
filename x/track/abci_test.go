package track

import (
	"fmt"
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestEndBlocker(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	// assure that play pool have rewards
	playPool := trackKeeper.GetPlayPool(ctx)
	fmt.Printf("Play pool: %s", playPool)
	fmt.Println()

	// Play tracks
	trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	trackKeeper.PlayTrack(ctx, addrDels[1], 1)
	trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	trackKeeper.PlayTrack(ctx, addrDels[0], 2)

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
	//fmt.Printf("%s", trackKeeper.GetAllPlays(ctx))

	playCollection := make(map[uint64]types.PlayReward)
	totalStreamPower := sdk.ZeroInt()

	trackKeeper.IterateAllPlays(ctx, func(play types.Play) bool {
		reward, ok := playCollection[play.TrackId]
		if !ok {
			reward = types.PlayReward{
				Streams:      play.Streams,
				Users:        sdk.NewInt(1),
				UsersPower:   play.Shares,
				StreamsPower: play.Streams.Quo(sdk.NewInt(1)).Mul(play.Shares.TruncateInt()),
				Reward:       sdk.NewInt64Coin("ubtsg", 0),
			}

		} else {
			reward.Streams = reward.Streams.Add(play.Streams)
			reward.Users = reward.Users.Add(sdk.NewInt(1))
			reward.UsersPower = reward.UsersPower.Add(play.Shares)
			reward.StreamsPower = reward.Streams.Quo(reward.Users).Mul(reward.UsersPower.TruncateInt())
		}
		playCollection[play.TrackId] = reward

		return false
	})

	// Calculate total Stream Power
	for play := range playCollection {
		totalStreamPower = totalStreamPower.Add(playCollection[play].StreamsPower)
		fmt.Printf("TSP: %s ", totalStreamPower)
		fmt.Printf("SP: %s", playCollection[play].StreamsPower)
		fmt.Println()
	}

	// Calculate price stream power
	priceStreamPower := playPool.Rewards.AmountOf("ubtsg").Quo(totalStreamPower)
	for play := range playCollection {
		item, ok := playCollection[play]
		if ok {
			item.Reward = item.Reward.Add(sdk.NewInt64Coin("ubtsg", priceStreamPower.Mul(item.StreamsPower).Int64()))
		}
		playCollection[play] = item

		fmt.Println()
		fmt.Printf("TrackID: %v ", play)
		fmt.Printf("Streams: %v ", playCollection[play].Streams)
		fmt.Printf("Users: %v ", playCollection[play].Users)
		fmt.Printf("UP: %v ", playCollection[play].UsersPower)
		fmt.Printf("SP: %v ", playCollection[play].StreamsPower)
		fmt.Printf("Reward: %v ", playCollection[play].Reward)
		fmt.Printf("priceStreamPower: %v ", priceStreamPower)
		fmt.Printf("totalStreamPower: %v ", totalStreamPower)
		fmt.Println()
	}

}
