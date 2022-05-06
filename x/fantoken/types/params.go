package types

import (
	"fmt"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyIssueFee = []byte("IssueFee")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIssueFee, &p.IssueFee, validateIssueFee),
	}
}

// NewParams constructs a new Params instance
func NewParams(issueFee sdk.Coin) Params {
	return Params{
		IssueFee: issueFee,
	}
}

// ParamKeyTable returns the TypeTable for the token module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams return the default params
func DefaultParams() Params {
	return Params{
		IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ValidateParams validates the given params
func ValidateParams(p Params) error {
	if err := validateIssueFee(p.IssueFee); err != nil {
		return err
	}

	return nil
}

func validateIssueFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("issue fee should not be negative")
	}
	return nil
}
