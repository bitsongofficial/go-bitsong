package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrLaunchPadDoesNotExist             = sdkerrors.Register(ModuleName, 2, "launchpad does not exist")
	ErrNotCollectionAuthority            = sdkerrors.Register(ModuleName, 3, "not authority of the collection")
	ErrNotLaunchPadAuthority             = sdkerrors.Register(ModuleName, 4, "not authority of the launchpad")
	ErrLaunchPadNotLiveTime              = sdkerrors.Register(ModuleName, 5, "launchpad is not live yet")
	ErrCannotExceedMaxMintParameter      = sdkerrors.Register(ModuleName, 6, "cannot exceed max mint parameter")
	ErrInsufficientMintableNftsRemaining = sdkerrors.Register(ModuleName, 7, "insufficient mintable nfts remaining")
)
