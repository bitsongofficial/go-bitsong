package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Artist struct {
	ID          ID             `json:"id"`
	Name        string         `json:"name"`
	URLs        URLs           `json:"urls"`
	Genres      []string       `json:"genres"`
	MetadataURI string         `json:"metadata"`
	Creator     sdk.AccAddress `json:"creator"`
}
