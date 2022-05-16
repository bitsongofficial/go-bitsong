package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrAuctionDoesNotExist = sdkerrors.Register(ModuleName, 2, "auction does not exist")
)
