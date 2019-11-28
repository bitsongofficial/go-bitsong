package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "artist"

	CodeInvalidGenesis    sdk.CodeType = 1
	CodeInvalidArtistName sdk.CodeType = 2
	CodeUnknownArtist     sdk.CodeType = 3
)

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrInvalidArtistName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistName, msg)
}

func ErrUnknownArtist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownArtist, msg)
}
