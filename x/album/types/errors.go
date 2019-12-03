package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidAlbumType                 sdk.CodeType = 0
	CodeInvalidAlbumName                 sdk.CodeType = 1
	CodeInvalidAlbumReleaseDate          sdk.CodeType = 2
	CodeInvalidAlbumReleaseDatePrecision sdk.CodeType = 3
	CodeInvalidGenesis                   sdk.CodeType = 4
	CodeUnknownAlbum                     sdk.CodeType = 5
	CodeInvalidAlbumStatus               sdk.CodeType = 6
)

func ErrInvalidAlbumType(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumType, msg)
}

func ErrInvalidAlbumName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumName, msg)
}

func ErrInvalidAlbumReleaseDate(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumReleaseDate, msg)
}

func ErrInvalidAlbumReleaseDatePrecision(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumReleaseDatePrecision, msg)
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrUnknownAlbum(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownAlbum, msg)
}

func ErrInvalidAlbumStatus(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumStatus, msg)
}
