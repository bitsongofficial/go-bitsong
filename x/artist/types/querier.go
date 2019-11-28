package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the Artist Querier
const (
	QueryParams  = "params"
	QueryArtists = "artists"
	QueryArtist  = "artist"
)

// Params for queries:
// - 'custom/artist/artist'
type QueryArtistParams struct {
	ArtistID uint64
}

// creates a new instance of QueryArtistParams
func NewQueryArtistParams(artistID uint64) QueryArtistParams {
	return QueryArtistParams{
		ArtistID: artistID,
	}
}

// Params for query 'custom/artist/artists'
type QueryArtistsParams struct {
	Owner        sdk.AccAddress
	ArtistStatus ArtistStatus
	Limit        uint64
}

// creates a new instance of QueryArtistsParams
func NewQueryArtistsParams(owner sdk.AccAddress, status ArtistStatus, limit uint64) QueryArtistsParams {
	return QueryArtistsParams{
		Owner:        owner,
		ArtistStatus: status,
		Limit:        limit,
	}
}
