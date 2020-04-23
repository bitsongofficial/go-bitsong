package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TrackTotals struct {
	Streams     uint64   `json:"streams" yaml:"streams"`
	Rewards   sdk.Coin `json:"rewards" yaml:"rewards"`
	Accounts uint64   `json:"accounts" yaml:"accounts"`
}

func (tt TrackTotals) String() string {
	return fmt.Sprintf(`Streams: %d - Accounts: %d - Rewards: %s`, tt.Streams, tt.Accounts, tt.Rewards.String())
}
