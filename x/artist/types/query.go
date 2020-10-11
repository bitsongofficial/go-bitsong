package types

const (
	QueryArtist = "artist"
)

type QueryArtistParams struct {
	ID string
}

func NewQueryArtistParams(id string) QueryArtistParams {
	return QueryArtistParams{ID: id}
}
