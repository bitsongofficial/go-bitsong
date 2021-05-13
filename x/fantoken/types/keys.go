package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the token module
	ModuleName = "fantoken"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the token module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the token module
	RouterKey string = ModuleName

	// DefaultParamspace is the default name for parameter store
	DefaultParamspace = ModuleName
)

var (
	// PrefixFanTokenForSymbol defines a symbol prefix for the fan token
	PrefixFanTokenForSymbol = []byte{0x01}
	// PrefixTokenForMinUint defines the min unit prefix for the token
	PrefixTokenForDenom = []byte{0x02}
	// PrefixFanTokens defines a prefix for the fan tokens
	PrefixFanTokens = []byte{0x03}
	// PeffxBurnFanTokenAmt defines a prefix for the amount of fan token burnt
	PefixBurnFanTokenAmt = []byte{0x04}
)

// KeySymbol returns the key of the token with the specified symbol
func KeySymbol(symbol string) []byte {
	return append(PrefixFanTokenForSymbol, []byte(symbol)...)
}

// KeyDenom returns the key of the token with the specified denom
func KeyDenom(denom string) []byte {
	return append(PrefixTokenForDenom, []byte(denom)...)
}

// KeyFanTokens returns the key of the specified owner and symbol. Intended for querying all fan tokens of an owner
func KeyFanTokens(owner sdk.AccAddress, symbol string) []byte {
	return append(append(PrefixFanTokens, owner.Bytes()...), []byte(symbol)...)
}

// KeyBurnTokenAmt returns the key of the specified min unit.
func KeyBurnFanTokenAmt(denom string) []byte {
	return append(PefixBurnFanTokenAmt, []byte(denom)...)
}
