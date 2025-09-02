package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	AddressCodec() address.Codec
}

type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
	HasSupply(ctx context.Context, denom string) bool

	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)

	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
