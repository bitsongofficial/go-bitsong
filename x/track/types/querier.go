package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the Track Querier
const (
	QueryParams = "params"
	QueryTracks = "tracks"
	QueryTrack  = "track"
)

// Params for queries:
// - 'custom/track/track'
type QueryTrackParams struct {
	TrackID uint64
}

// creates a new instance of QueryTracParams
func NewQueryTrackParams(trackID uint64) QueryTrackParams {
	return QueryTrackParams{
		TrackID: trackID,
	}
}

// Params for query 'custom/track/tracks'
type QueryTracksParams struct {
	Owner       sdk.AccAddress
	TrackStatus TrackStatus
	Limit       uint64
}

// creates a new instance of QueryTracksParams
func NewQueryTracksParams(owner sdk.AccAddress, status TrackStatus, limit uint64) QueryTracksParams {
	return QueryTracksParams{
		Owner:       owner,
		TrackStatus: status,
		Limit:       limit,
	}
}
