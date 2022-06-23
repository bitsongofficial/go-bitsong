package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MinimumSymbolLen is the minimum limitation for the length of the fantoken's symbol
	MinimumSymbolLen = 1
	// MaximumSymbolLen is the maximum limitation for the length of the fantoken's symbol
	MaximumSymbolLen = 64
	// MinimumNameLen is the minimum limitation for the length of the fantoken's name
	MinimumNameLen = 0
	// MaximumNameLen is the maximum limitation for the length of the fantoken's name
	MaximumNameLen = 128
	// MinimumUriLen is the minimum limitation for the length of the fantoken's uri
	MinimumUriLen = 0
	// MaximumUriLen is the maximum limitation for the length of the fantoken's uri
	MaximumUriLen = 512
)

var (
	regexpSymbolFmt = fmt.Sprintf("^[a-z0-9]{%d,%d}$", MinimumSymbolLen-1, MaximumSymbolLen-1)
	regexpSymbol    = regexp.MustCompile(regexpSymbolFmt).MatchString
)

// ValidateDenom checks if the given denom is valid
func ValidateDenom(denom string) error {
	if !strings.HasPrefix(denom, "ft") {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom: %s, denom starts with ft", denom)
	}

	return sdk.ValidateDenom(denom)
}

// ValidateName verifies whether the given name is valid
func ValidateName(name string) error {
	if len(strings.TrimSpace(name)) > MaximumNameLen {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid fantoken name %s, only accepts length (%d, %d]", name, MinimumNameLen, MaximumNameLen)
	}

	return nil
}

// ValidateSymbol checks if the given symbol is valid
func ValidateSymbol(symbol string) error {
	if !regexpSymbol(strings.TrimSpace(symbol)) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, only accepts english lowercase letters and numbers, length [%d, %d], and begin with an english letter, regexp: %s", symbol, MinimumSymbolLen, MaximumSymbolLen, regexpSymbolFmt)
	}

	return nil
}

// ValidateAmount checks if the given amount is positive amount
func ValidateAmount(amount sdk.Int) error {
	if amount.IsZero() || amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidMaxSupply, "invalid fantoken amount %d, only accepts positive amount", amount)
	}
	return nil
}

// ValidateUri checks if the given uri is valid
func ValidateUri(uri string) error {
	if len(strings.TrimSpace(uri)) > MaximumUriLen {
		return sdkerrors.Wrapf(ErrInvalidUri, "invalid uri: %s, uri only accepts length (%d, %d]", uri, MinimumUriLen, MaximumUriLen)
	}

	return nil
}

func ValidateFees(issueFee, mintFee, burnFee sdk.Coin) error {
	if err := issueFee.Validate(); err != nil {
		return err
	}

	if err := mintFee.Validate(); err != nil {
		return err
	}

	if err := burnFee.Validate(); err != nil {
		return err
	}

	return nil
}
