package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/reward/types"
)

const (
	DefaultParamspace = types.ModuleName
)

var (
	RewardPoolKey          = []byte{0x00}
	ParamStoreKeyRewardTax = []byte("rewardtax")
)
