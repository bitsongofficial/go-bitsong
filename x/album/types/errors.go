package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidAlbumType        sdk.CodeType = 0
	CodeInvalidAlbumName        sdk.CodeType = 1
	CodeInvalidGenesis          sdk.CodeType = 2
	CodeUnknownAlbum            sdk.CodeType = 3
	CodeInvalidAlbumStatus      sdk.CodeType = 4
	CodeInvalidAlbumMetadataUri sdk.CodeType = 5
	CodeUnknownTrack            sdk.CodeType = 20
)

func ErrInvalidAlbumType(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumType, msg)
}

func ErrInvalidAlbumName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumName, msg)
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

func ErrUnknownTrack(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownTrack, msg)
}

func ErrInvalidAlbumMetadataURI(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAlbumMetadataUri, msg)
}
