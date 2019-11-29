package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis           sdk.CodeType = 1
	CodeInvalidArtistName        sdk.CodeType = 2
	CodeUnknownArtist            sdk.CodeType = 3
	CodeInvalidArtistImageHeight sdk.CodeType = 4
	CodeInvalidArtistImageWidth  sdk.CodeType = 5
	CodeInvalidArtistImageCid    sdk.CodeType = 6
	CodeUnknownOwner             sdk.CodeType = 7
	CodeInvalidArtistStatus      sdk.CodeType = 8
	CodeUnknownModerator         sdk.CodeType = 9
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

func ErrInvalidArtistImageHeight(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistImageHeight, msg)
}

func ErrInvalidArtistImageWidth(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistImageWidth, msg)
}

func ErrInvalidArtistImageCid(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistImageCid, msg)
}

func ErrInvalidArtistStatus(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistStatus, msg)
}

func ErrUnknownOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownOwner, msg)
}

func ErrUnknownModerator(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownModerator, msg)
}
