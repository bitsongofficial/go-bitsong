package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines the default reward module parameter subspace
const DefaultParamspace = ModuleName

// Default parameter values
var (
	ParamStoreKeyRewardTax = []byte("rewardtax")

	DefaultRewardTx = sdk.NewDecWithPrec(1, 2)
)

var _ paramtypes.ParamSet = &Params{}

// Params defines the parameters for the reward module.
type Params struct {
	ParamsStoreKeyRewardTax sdk.Dec
}

// ParamKeyTable for reward module
func ParamKeyTable() params.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of reward module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyRewardTax, &p.ParamsStoreKeyRewardTax, validateRewardTax),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		ParamsStoreKeyRewardTax: DefaultRewardTx,
	}
}

func validateRewardTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("invalid reward tax amount: %s", v)
	}

	// TODO: Add other validations

	return nil
}
