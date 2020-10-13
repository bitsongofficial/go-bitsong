package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) exported.Account
	SetAccount(sdk.Context, exported.Account)
	IterateAccounts(ctx sdk.Context, process func(exported.Account) (stop bool))
}
