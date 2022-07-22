package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrCandyMachineDoesNotExist = sdkerrors.Register(ModuleName, 2, "candy machine does not exist")
)
