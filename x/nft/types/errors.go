package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrCollectionAlreadyExists = sdkerrors.Register(ModuleName, 1, "invalid collection: already exists")
	ErrCollectionNotFound      = sdkerrors.Register(ModuleName, 2, "invalid collection: not found")
)
