package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "artist"

	CodeUnknownArtist          sdk.CodeType = 1
	CodeAlreadyVerifiedArtist  sdk.CodeType = 2
	CodeInvalidGenesis         sdk.CodeType = 3
	CodeInvalidArtistStatus    sdk.CodeType = 4
	CodeInvalidArtistMeta      sdk.CodeType = 5
	CodeArtistHandlerNotExists sdk.CodeType = 6
)

func ErrUnknownArtist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownArtist, msg)
}

func ErrAlreadyVerifiedArtist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyVerifiedArtist, msg)
}

func ErrInvalidArtistStatus(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistStatus, msg)
}

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrInvalidArtistMeta(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidArtistMeta, msg)
}

func ErrNoArtistHandlerExists(codespace sdk.CodespaceType, meta interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeArtistHandlerNotExists, fmt.Sprintf("'%T' does not have a corresponding handler", meta))
}
