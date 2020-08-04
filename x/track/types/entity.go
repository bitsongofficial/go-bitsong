package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Entity struct {
	Shares  sdk.Int        `json:"share" yaml:"share"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}
