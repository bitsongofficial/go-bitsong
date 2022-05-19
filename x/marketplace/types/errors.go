package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrAuctionDoesNotExist   = sdkerrors.Register(ModuleName, 2, "auction does not exist")
	ErrNotAuctionAuthority   = sdkerrors.Register(ModuleName, 3, "not an auction authority")
	ErrAuctionAlreadyStarted = sdkerrors.Register(ModuleName, 4, "auction already started")
	ErrAuctionAlreadyEnded   = sdkerrors.Register(ModuleName, 5, "auction already ended")
)
