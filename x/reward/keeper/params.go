package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

func ParamKeyTable() params.KeyTable {
	keyTable := params.KeyTable{}
	keyTable.RegisterType(params.ParamSetPair{
		Key:   ParamStoreKeyRewardTax,
		Value: sdk.Dec{},
	})
	return keyTable
}

func (k Keeper) GetRewardTax(ctx sdk.Context) sdk.Dec {
	var percent sdk.Dec
	k.paramSpace.Get(ctx, ParamStoreKeyRewardTax, &percent)
	return percent
}

func (k Keeper) SetRewardTax(ctx sdk.Context, percent sdk.Dec) {
	k.paramSpace.Set(ctx, ParamStoreKeyRewardTax, &percent)
}
