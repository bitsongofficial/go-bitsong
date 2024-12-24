package types

import (
	"fmt"

	"cosmossdk.io/math"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyIssueFee = []byte("IssueFee")
	KeyMintFee  = []byte("MintFee")
	KeyBurnFee  = []byte("BurnFee")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIssueFee, &p.IssueFee, validateFee),
		paramtypes.NewParamSetPair(KeyMintFee, &p.MintFee, validateFee),
		paramtypes.NewParamSetPair(KeyBurnFee, &p.BurnFee, validateFee),
	}
}

// NewParams constructs a new Params instance
func NewParams(issueFee, mintFee, burnFee, transferFee sdk.Coin) Params {
	return Params{
		IssueFee: issueFee,
		MintFee:  mintFee,
		BurnFee:  burnFee,
	}
}

// ParamKeyTable returns the TypeTable for the token module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams return the default params
func DefaultParams() Params {
	return Params{
		IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000_000)),
		MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
		BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates the given params
func (p Params) Validate() error {
	if err := validateFee(p.IssueFee); err != nil {
		return err
	}

	if err := validateFee(p.MintFee); err != nil {
		return err
	}

	if err := validateFee(p.BurnFee); err != nil {
		return err
	}

	return nil
}

func validateFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("fee should not be negative")
	}
	return nil
}
