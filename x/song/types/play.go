package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Play struct {
	AccountAddress sdk.AccAddress `json:"account_address"`
	SongId         sdk.Int        `json:"song_id"`
	Shares         sdk.Dec        `json:"shares"`
	Streams        sdk.Int        `json:"streams"`
}
