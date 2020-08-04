package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Share struct {
	TrackID string         `json:"track_id" yaml:"track_id"`
	Entity  sdk.AccAddress `json:"entity" yaml:"entity"`
	Shares  sdk.Coin       `json:"shares" yaml:"shares"`
}
