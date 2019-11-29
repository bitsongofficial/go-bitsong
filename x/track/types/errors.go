package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis    sdk.CodeType = 0
	CodeInvalidTrackTitle sdk.CodeType = 1
)

func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrInvalidTrackTitle(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTrackTitle, msg)
}
