package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis          sdk.CodeType = 0
	CodeInvalidTrackTitle       sdk.CodeType = 1
	CodeUnknownTrack            sdk.CodeType = 2
	CodeInvalidTrackStatus      sdk.CodeType = 3
	CodeInvalidTrackMetadataURI sdk.CodeType = 4
)

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrInvalidTrackTitle(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrackTitle, msg)
}

func ErrUnknownTrack(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownTrack, msg)
}

func ErrInvalidTrackStatus(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrackStatus, msg)
}

func ErrInvalidTrackMetadataURI(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrackMetadataURI, msg)
}
