package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis           sdk.CodeType = 1
	CodeInvalidArtistName        sdk.CodeType = 2
	CodeInvalidArtistMetadataURI sdk.CodeType = 3
	CodeInvalidArtistStatus      sdk.CodeType = 4
	CodeUnknownArtist            sdk.CodeType = 5
	CodeUnknownOwner             sdk.CodeType = 6
)

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrInvalidArtistName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistName, msg)
}

func ErrInvalidArtistMetadataURI(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistMetadataURI, msg)
}

func ErrUnknownArtist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownArtist, msg)
}

func ErrInvalidArtistStatus(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistStatus, msg)
}

func ErrUnknownOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownOwner, msg)
}
