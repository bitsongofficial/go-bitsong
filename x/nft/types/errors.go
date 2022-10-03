package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrNFTDoesNotExist                = sdkerrors.Register(ModuleName, 2, "NFT does not exist")
	ErrMetadataDoesNotExist           = sdkerrors.Register(ModuleName, 3, "metadata does not exist")
	ErrCollectionDoesNotExist         = sdkerrors.Register(ModuleName, 4, "collection does not exist")
	ErrNotNFTOwner                    = sdkerrors.Register(ModuleName, 5, "not the owner of the NFT")
	ErrNotEnoughPermission            = sdkerrors.Register(ModuleName, 6, "not enough permission")
	ErrMetadataImmutable              = sdkerrors.Register(ModuleName, 7, "metadata is immutable")
	ErrPrimarySaleAlreadyHappened     = sdkerrors.Register(ModuleName, 8, "primary sale already happened")
	ErrInvalidSellerFeeBasisPoints    = sdkerrors.Register(ModuleName, 9, "invalid seller fee basis points")
	ErrNotMasterEditionNft            = sdkerrors.Register(ModuleName, 10, "not master edition NFT")
	ErrAlreadyReachedEditionMaxSupply = sdkerrors.Register(ModuleName, 11, "already reached edition maximum supply")
	ErrInvalidNftId                   = sdkerrors.Register(ModuleName, 12, "invalid NFT id")
	ErrMasterEditionNftDoesNotExists  = sdkerrors.Register(ModuleName, 13, "master edition NFT does not exist")
)
