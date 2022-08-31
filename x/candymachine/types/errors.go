package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrCandyMachineDoesNotExist          = sdkerrors.Register(ModuleName, 2, "candy machine does not exist")
	ErrNotCollectionAuthority            = sdkerrors.Register(ModuleName, 3, "not authority of the collection")
	ErrNotCandyMachineAuthority          = sdkerrors.Register(ModuleName, 4, "not authority of the candy machine")
	ErrCandyMachineNotLiveTime           = sdkerrors.Register(ModuleName, 5, "candy machine is not live yet")
	ErrCannotExceedMaxMintParameter      = sdkerrors.Register(ModuleName, 6, "cannot exceed max mint parameter")
	ErrInsufficientMintableNftsRemaining = sdkerrors.Register(ModuleName, 7, "insufficient mintable nfts remaining")
)
