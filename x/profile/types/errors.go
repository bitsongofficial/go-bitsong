package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ErrCodeProfileNotFound    = 1
	ErrCodeProfileCreateError = 2
)

var (
	DefaultCodespace = ModuleName

	ErrProfileNotFound    = sdkerrors.Register(ModuleName, ErrCodeProfileNotFound, "profile not found")
	ErrProfileCreateError = sdkerrors.Register(ModuleName, ErrCodeProfileCreateError, "profile create error")
)
