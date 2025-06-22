package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type DistrKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	GetParamSet(ctx sdk.Context, ps paramstypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramstypes.ParamSet)
	HasKeyTable() bool
	WithKeyTable(table paramstypes.KeyTable) paramstypes.Subspace
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetSupply(ctx context.Context, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	//GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	//SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	//SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
