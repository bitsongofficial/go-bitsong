package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ErrCodeReleaseNotFound    = 1
	ErrCodeReleaseCreateError = 2
)

var (
	DefaultCodespace = ModuleName

	ErrReleaseNotFound    = sdkerrors.Register(ModuleName, ErrCodeReleaseNotFound, "release not found")
	ErrReleaseCreateError = sdkerrors.Register(ModuleName, ErrCodeReleaseCreateError, "release create error")
)
