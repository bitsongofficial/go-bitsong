package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyCreationFee = []byte("CreationFee")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCreationFee, &p.CreationFee, validateCreationFee),
	}
}

// NewParams constructs a new Params instance
func NewParams(creationFee sdk.Coin) Params {
	return Params{
		CreationFee: creationFee,
	}
}

// ParamKeyTable returns the TypeTable for the module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams return the default params
func DefaultParams() Params {
	return Params{
		CreationFee: sdk.NewInt64Coin(sdk.DefaultBondDenom, 1_000_000),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates the given params
func (p Params) Validate() error {
	if err := validateCreationFee(p.CreationFee); err != nil {
		return err
	}

	return nil
}

func validateCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid creation fee: %+v", i)
	}

	return nil
}
