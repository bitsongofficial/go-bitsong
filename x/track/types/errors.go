package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	DefaultCodespace = ModuleName
	ErrUnknownTrack  = sdkerrors.Register(ModuleName, 1, "unknown track")
	ErrUnknownShare  = sdkerrors.Register(ModuleName, 2, "unknown share")
	ErrInvalidAmount = sdkerrors.Register(ModuleName, 3, "invalid amount")
)
