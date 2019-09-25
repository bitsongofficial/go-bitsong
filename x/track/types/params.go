package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	KeyPlayTax  = []byte("PlayTax")
	PlayPoolKey = []byte("PlayPool")
)

type Params struct {
	PlayTax sdk.Dec `json:"play_tax"`
}

// ParamTable for token module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(playTax sdk.Dec, startingTrackID uint64) Params {
	return Params{
		PlayTax: playTax,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		PlayTax: sdk.NewDecWithPrec(30, 2),
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyPlayTax, Value: &p.PlayTax},
	}
}
