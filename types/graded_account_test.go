package types

import (
	"testing"
	"time"

	"github.com/BitSongOfficial/go-bitsong/types/assets"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
)

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

var (
	angelSchedule        VestingSchedule
	seedSchedule         VestingSchedule
	privateSchedule      VestingSchedule
	privateBonusSchedule VestingSchedule
	employeeSchedule     VestingSchedule
	timeGenesisString    = "2019-04-23 23:00:00 -0800 PST"
	monthlyTimes         []int64
	timeGenesis          time.Time
)

// initialize the times!
func init() {
	var err error
	timeLayoutString := "2006-01-02 15:04:05 -0700 MST"
	timeGenesis, err = time.Parse(timeLayoutString, timeGenesisString)
	if err != nil {
		panic(err)
	}

	monthlyTimes = []int64{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 12; j++ {
			monthlyTimes = append(monthlyTimes, timeGenesis.AddDate(i, j, 0).Unix())
		}
	}

	angelSchedule = NewVestingSchedule(assets.MicroBitSongDenom, []Schedule{
		NewSchedule(monthlyTimes[1], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[2], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[3], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[12], sdk.NewDecWithPrec(70, 2)),
	})

	seedSchedule = NewVestingSchedule(assets.MicroBitSongDenom, []Schedule{
		NewSchedule(monthlyTimes[1], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[2], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[3], sdk.NewDecWithPrec(10, 2)),
		NewSchedule(monthlyTimes[10], sdk.NewDecWithPrec(70, 2)),
	})

	privateSchedule = NewVestingSchedule(assets.MicroBitSongDenom, []Schedule{
		NewSchedule(monthlyTimes[3], sdk.NewDecWithPrec(16, 2)),
		NewSchedule(monthlyTimes[4], sdk.NewDecWithPrec(17, 2)),
		NewSchedule(monthlyTimes[5], sdk.NewDecWithPrec(16, 2)),
		NewSchedule(monthlyTimes[6], sdk.NewDecWithPrec(17, 2)),
		NewSchedule(monthlyTimes[7], sdk.NewDecWithPrec(17, 2)),
		NewSchedule(monthlyTimes[8], sdk.NewDecWithPrec(17, 2)),
	})

	privateBonusSchedule = NewVestingSchedule(assets.MicroBitSongDenom, []Schedule{
		NewSchedule(monthlyTimes[6], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[7], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[8], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[9], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[10], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[11], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[12], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[13], sdk.NewDecWithPrec(8, 2)),
		NewSchedule(monthlyTimes[14], sdk.NewDecWithPrec(9, 2)),
		NewSchedule(monthlyTimes[15], sdk.NewDecWithPrec(9, 2)),
		NewSchedule(monthlyTimes[16], sdk.NewDecWithPrec(9, 2)),
		NewSchedule(monthlyTimes[17], sdk.NewDecWithPrec(9, 2)),
	})

	employeeSchedule = NewVestingSchedule(assets.MicroBitSongDenom, []Schedule{
		NewSchedule(monthlyTimes[0], sdk.NewDecWithPrec(5, 2)),
		NewSchedule(monthlyTimes[12], sdk.NewDecWithPrec(29, 2)),
		NewSchedule(monthlyTimes[24], sdk.NewDecWithPrec(33, 2)),
		NewSchedule(monthlyTimes[36], sdk.NewDecWithPrec(33, 2)),
	})

}

func scaleCoins(scale float64, denom string, input sdk.Coins) sdk.Coins {
	output := sdk.Coins{}
	for _, coin := range input {
		if coin.Denom != denom {
			continue
		}

		decScale := sdk.NewDecWithPrec(int64(scale*100), 2)
		output = append(output, sdk.NewCoin(coin.Denom, decScale.MulInt(coin.Amount).RoundInt()))
	}
	return output
}