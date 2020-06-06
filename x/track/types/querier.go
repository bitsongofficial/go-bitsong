package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams        = "params"
	QueryID            = "id"
	QueryCreatorTracks = "creator_tracks"
)

type QueryCreatorTracksParams struct {
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
}

// Params for queries
type QueryTrackParams struct {
	ID uint64 `json:"id" yaml:"id"`
}

// creates a new instance of QueryContentParams
func NewQueryContentParams(id uint64) QueryTrackParams {
	return QueryTrackParams{
		ID: id,
	}
}
