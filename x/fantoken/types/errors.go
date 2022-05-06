//nolint
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// fantoken module errors
var (
	ErrInvalidName        = sdkerrors.Register(ModuleName, 1, "invalid fantoken name")
	ErrInvalidDenom       = sdkerrors.Register(ModuleName, 2, "invalid fantoken denom")
	ErrInvalidSymbol      = sdkerrors.Register(ModuleName, 3, "invalid standard symbol")
	ErrInvalidMaxSupply   = sdkerrors.Register(ModuleName, 4, "invalid fantoken maximum supply")
	ErrDenomAlreadyExists = sdkerrors.Register(ModuleName, 5, "denom already exists")
	ErrFanTokenNotExists  = sdkerrors.Register(ModuleName, 6, "fantoken does not exist")
	ErrInvalidToAddress   = sdkerrors.Register(ModuleName, 7, "the new owner must not be same as the original owner")
	ErrInvalidOwner       = sdkerrors.Register(ModuleName, 8, "invalid fantoken owner")
	ErrNotMintable        = sdkerrors.Register(ModuleName, 9, "fantoken is not mintable")
	ErrNotFoundTokenAmt   = sdkerrors.Register(ModuleName, 10, "burned fantoken amount not found")
	ErrInvalidAmount      = sdkerrors.Register(ModuleName, 11, "invalid amount")
	ErrInvalidUri         = sdkerrors.Register(ModuleName, 12, "invalid uri length")
)
