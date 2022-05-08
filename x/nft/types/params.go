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
	KeyIssuePrice = []byte("IssuePrice")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIssuePrice, &p.IssuePrice, validateIssuePrice),
	}
}

// NewParams constructs a new Params instance
func NewParams(issuePrice sdk.Coin) Params {
	return Params{
		IssuePrice: issuePrice,
	}
}

// ParamKeyTable returns the TypeTable for the module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams return the default params
func DefaultParams() Params {
	return Params{}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ValidateParams validates the given params
func ValidateParams(p Params) error {
	if err := validateIssuePrice(p.IssuePrice); err != nil {
		return err
	}

	return nil
}

func validateIssuePrice(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("issue price should not be negative")
	}
	return nil
}
