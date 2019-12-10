package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidDistributorTitle sdk.CodeType = 0
	CodeInvalidDistributor      sdk.CodeType = 1
)

func ErrInvalidDistributorName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDistributorTitle, msg)
}

func ErrInvalidDistributor(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDistributor, msg)
}
