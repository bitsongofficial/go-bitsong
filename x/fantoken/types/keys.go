package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the token module
	ModuleName = "fantoken"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// RouterKey is the msg router key for the token module
	RouterKey string = ModuleName
)

var (
	// PrefixFanTokenForDenom defines a denom prefix for the fan token
	PrefixFanTokenForDenom = []byte{0x01}

	// PrefixFanTokens defines a prefix for the fan tokens
	PrefixFanTokens = []byte{0x02}
)

// KeyDenom returns the key of the token with the specified denom
func KeyDenom(denom string) []byte {
	return append(PrefixFanTokenForDenom, []byte(denom)...)
}

// KeyFanTokens returns the key of the specified owner and denom. Intended for querying all fan tokens of an owner
func KeyFanTokens(owner sdk.AccAddress, denom string) []byte {
	return append(append(PrefixFanTokens, owner.Bytes()...), []byte(denom)...)
}
