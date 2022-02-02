package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MinimumSymbolLen is the minimum limitation for the length of the token's symbol
	MinimumSymbolLen = 3
	// MaximumSymbolLen is the maximum limitation for the length of the token's symbol
	MaximumSymbolLen = 64
	// MaximumNameLen is the maximum limitation for the length of the token's name
	MaximumNameLen = 32
	// MinimumDenomLen is the minimum limitation for the length of the token's denom
	MinimumDenomLen = 3
	// MaximumMinUnitLen is the maximum limitation for the length of the token's denom
	MaximumDenomLen = 64
)

const (
	ReservedPeg  = "peg"
	ReservedIBC  = "ibc"
	ReservedSwap = "swap"
	ReservedHTLT = "htlt"
)

var (
	keywords          = strings.Join([]string{ReservedPeg, ReservedIBC, ReservedSwap, ReservedHTLT}, "|")
	regexpKeywordsFmt = fmt.Sprintf("^(%s).*", keywords)
	regexpKeyword     = regexp.MustCompile(regexpKeywordsFmt).MatchString

	regexpSymbolFmt = fmt.Sprintf("^[a-z][a-z0-9]{%d,%d}$", MinimumSymbolLen-1, MaximumSymbolLen-1)
	regexpSymbol    = regexp.MustCompile(regexpSymbolFmt).MatchString
)

// ValidateToken checks if the given token is valid
func ValidateToken(token FanToken) error {
	if len(token.Owner) > 0 {
		if _, err := sdk.AccAddressFromBech32(token.Owner); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}
	if err := ValidateName(token.Name); err != nil {
		return err
	}
	if err := ValidateSymbol(token.GetSymbol()); err != nil {
		return err
	}
	if err := ValidateDenom(token.GetDenom()); err != nil {
		return err
	}
	return nil
}

// ValidateDenom checks if the given denom is valid
func ValidateDenom(denom string) error {
	if !strings.HasPrefix(denom, "ft") {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom: %s, denom starts with ft", denom)
	}
	return ValidateKeywords(denom)
}

// ValidateName verifies whether the given name is legal
func ValidateName(name string) error {
	if len(name) == 0 || len(name) > MaximumNameLen {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid token name %s, only accepts length (0, %d]", name, MaximumNameLen)
	}
	return nil
}

// ValidateSymbol checks if the given symbol is valid
func ValidateSymbol(symbol string) error {
	if !regexpSymbol(symbol) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, only accepts english lowercase letters and numbers, length [%d, %d], and begin with an english letter, regexp: %s", symbol, MinimumSymbolLen, MaximumSymbolLen, regexpSymbolFmt)
	}
	return ValidateKeywords(symbol)
}

// ValidateKeywords checks if the given symbol begins with `TokenKeywords`
func ValidateKeywords(denom string) error {
	if regexpKeyword(denom) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid token: %s, can not begin with keyword: (%s)", denom, keywords)
	}
	return nil
}

// ValidateAmount checks if the given amount is positive amount
func ValidateAmount(amount sdk.Int) error {
	if amount.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidMaxSupply, "invalid token amount %d, only accepts positive amount", amount)
	}
	return nil
}
