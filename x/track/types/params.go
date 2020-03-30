package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines the default track module parameter subspace
const DefaultParamspace = ModuleName

// Parameter store key
var (
	ParamStoreKeyDepositParams = []byte("trackdepositparams")

	DefaultDepositParams = DepositParams{}
)

var _ paramtypes.ParamSet = &Params{}

// Params defines the parameters for the track module.
type Params struct {
	DepositParams DepositParams
}

// Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of reward module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyDepositParams, &p.DepositParams, validateDepositParams),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		DepositParams: DefaultDepositParams,
	}
}

func validateDepositParams(i interface{}) error {
	v, ok := i.(DepositParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MinDeposit.IsAnyNegative() {
		return fmt.Errorf("invalid min deposit value: %s", v)
	}

	if v.MaxDepositPeriod.Microseconds() == 0 {
		return fmt.Errorf("invalid max deposit period: %s", v.MaxDepositPeriod)
	}

	// TODO: Add other validations

	return nil
}

// Param around deposits for artists
type DepositParams struct {
	MinDeposit       sdk.Coins     `json:"min_deposit,omitempty" yaml:"min_deposit,omitempty"`               //  Minimum deposit for the track to be self-verified.
	MaxDepositPeriod time.Duration `json:"max_deposit_period,omitempty" yaml:"max_deposit_period,omitempty"` //  Maximum period for Btsg holders to deposit on a specific artist.
}

// NewDepositParams creates a new DepositParams object
func NewDepositParams(minDeposit sdk.Coins, maxDepositPeriod time.Duration) DepositParams {
	return DepositParams{
		MinDeposit:       minDeposit,
		MaxDepositPeriod: maxDepositPeriod,
	}
}

func (dp DepositParams) String() string {
	return fmt.Sprintf(`Deposit Params:
  Min Deposit:        %s
  Max Deposit Period: %s`, dp.MinDeposit, dp.MaxDepositPeriod)
}

// Checks equality of DepositParams
func (dp DepositParams) Equal(dp2 DepositParams) bool {
	return dp.MinDeposit.IsEqual(dp2.MinDeposit) && dp.MaxDepositPeriod == dp2.MaxDepositPeriod
}
