package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	RewardPoolKey = []byte{0x00}
)

func (k Keeper) GetRewardTax(ctx sdk.Context) sdk.Dec {
	var percent sdk.Dec
	k.paramSpace.Get(ctx, types.ParamStoreKeyRewardTax, &percent)
	return percent
}

func (k Keeper) SetRewardTax(ctx sdk.Context, percent sdk.Dec) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyRewardTax, &percent)
}
