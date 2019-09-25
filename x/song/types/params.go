package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	KeySongTax = []byte("SongTax")
)

type Params struct {
	SongTax sdk.Dec `json:"song_tax"`
}

// ParamTable for token module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(songTax sdk.Dec) Params {
	return Params{
		SongTax: songTax,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		SongTax: sdk.NewDecWithPrec(30, 2),
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeySongTax, Value: &p.SongTax},
	}
}
