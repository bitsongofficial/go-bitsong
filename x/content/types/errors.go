package types

import (
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	DefaultCodespace = ModuleName
	ErrCreateFailed  = sdkErrors.Register(DefaultCodespace, 1, "create hls contract failed")
)
