package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrNFTDoesNotExist        = sdkerrors.Register(ModuleName, 2, "nft does not exist")
	ErrMetadataDoesNotExist   = sdkerrors.Register(ModuleName, 3, "metadata does not exist")
	ErrCollectionDoesNotExist = sdkerrors.Register(ModuleName, 4, "collection does not exist")
)
