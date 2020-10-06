package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ErrCodeChannelNotFound    = 1
	ErrCodeChannelCreateError = 2
)

var (
	DefaultCodespace = ModuleName

	ErrChannelNotFound    = sdkerrors.Register(ModuleName, ErrCodeChannelNotFound, "channel not found")
	ErrChannelCreateError = sdkerrors.Register(ModuleName, ErrCodeChannelCreateError, "channel create error")
)
