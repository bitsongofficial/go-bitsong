package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Channel struct {
	Username  string         `json:"username" yaml:"username"`
	Hash      string         `json:"hash" yaml:"hash"`
	Owner     sdk.AccAddress `json:"owner" yaml:"owner"`
	CreatedAt time.Time      `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" yaml:"updated_at"`
}
