package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

// Parameter store key
var (
	ParamStoreKeyDepositParams = []byte("depositparams")
)

// Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyDepositParams, DepositParams{},
	)
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
