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
	KeyLaunchPadCreationPrice = []byte("LaunchPadCreationPrice")
	KeyLaunchpadMaxMint       = []byte("LaunchPadMaxMint")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLaunchPadCreationPrice, &p.LaunchpadCreationPrice, validateLaunchPadCreationPrice),
		paramtypes.NewParamSetPair(KeyLaunchpadMaxMint, &p.LaunchpadMaxMint, validateLaunchPadMaxMint),
	}
}

// NewParams constructs a new Params instance
func NewParams(launchpadCreationPrice sdk.Coin, launchpadMaxMint uint64) Params {
	return Params{
		LaunchpadCreationPrice: launchpadCreationPrice,
		LaunchpadMaxMint:       launchpadMaxMint,
	}
}

// ParamKeyTable returns the TypeTable for the module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams return the default params
func DefaultParams() Params {
	return Params{
		LaunchpadCreationPrice: sdk.NewInt64Coin("ubtsg", 0),
		LaunchpadMaxMint:       10000,
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ValidateParams validates the given params
func ValidateParams(p Params) error {
	if err := validateLaunchPadCreationPrice(p.LaunchpadCreationPrice); err != nil {
		return err
	}

	return nil
}

func validateLaunchPadCreationPrice(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("issue price should not be negative")
	}
	return nil
}

func validateLaunchPadMaxMint(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max mint should be positive value")
	}
	return nil
}
