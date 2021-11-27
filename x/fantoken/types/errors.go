//nolint
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// token module sentinel errors
var (
	ErrInvalidName         = sdkerrors.Register(ModuleName, 2, "invalid token name")
	ErrInvalidDenom        = sdkerrors.Register(ModuleName, 3, "invalid token denom")
	ErrInvalidSymbol       = sdkerrors.Register(ModuleName, 4, "invalid standard symbol")
	ErrInvalidInitSupply   = sdkerrors.Register(ModuleName, 5, "invalid token initial supply")
	ErrInvalidMaxSupply    = sdkerrors.Register(ModuleName, 6, "invalid token maximum supply")
	ErrInvalidScale        = sdkerrors.Register(ModuleName, 7, "invalid token scale")
	ErrSymbolAlreadyExists = sdkerrors.Register(ModuleName, 8, "symbol already exists")
	ErrDenomAlreadyExists  = sdkerrors.Register(ModuleName, 9, "denom already exists")
	ErrFanTokenNotExists   = sdkerrors.Register(ModuleName, 10, "fantoken does not exist")
	ErrInvalidToAddress    = sdkerrors.Register(ModuleName, 11, "the new owner must not be same as the original owner")
	ErrInvalidOwner        = sdkerrors.Register(ModuleName, 12, "invalid token owner")
	ErrNotMintable         = sdkerrors.Register(ModuleName, 13, "token is not mintable")
	ErrNotFoundTokenAmt    = sdkerrors.Register(ModuleName, 14, "burned token amount not found")
	ErrInvalidAmount       = sdkerrors.Register(ModuleName, 15, "invalid amount")
	ErrLessIssueFee        = sdkerrors.Register(ModuleName, 16, "less issue fee")
)
