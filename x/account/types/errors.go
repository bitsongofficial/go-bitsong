package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ErrCodeAccountNotFound    = 1
	ErrCodeAccountCreateError = 2
)

var (
	DefaultCodespace = ModuleName

	ErrAccountNotFound    = sdkerrors.Register(ModuleName, ErrCodeAccountNotFound, "account not found")
	ErrAccountCreateError = sdkerrors.Register(ModuleName, ErrCodeAccountCreateError, "account create error")
)
