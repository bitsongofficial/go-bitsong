// nolint
package types

import (
	sdkerrors "cosmossdk.io/errors"
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
	ErrInvalidAuthority   = sdkerrors.Register(ModuleName, 8, "invalid fantoken authority")
	ErrInvalidMinter      = sdkerrors.Register(ModuleName, 9, "invalid fantoken minter")
	ErrInvalidRecipient   = sdkerrors.Register(ModuleName, 10, "invalid fantoken recipient")
	ErrInvalidOwner       = sdkerrors.Register(ModuleName, 11, "the owner is empty or invalid")
	ErrNotFoundTokenAmt   = sdkerrors.Register(ModuleName, 12, "burned fantoken amount not found")
	ErrInvalidAmount      = sdkerrors.Register(ModuleName, 13, "invalid amount")
	ErrInvalidUri         = sdkerrors.Register(ModuleName, 14, "invalid uri length")
)
