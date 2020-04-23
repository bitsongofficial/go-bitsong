package types

// query endpoints supported by the Track Querier
const (
	QueryParams = "params"
	QueryTrack  = "track"
)

// Params for queries
type QueryTrackParams struct {
	TrackAddr string `json:"track_addr" yaml:"track_addr"`
}

// creates a new instance of QueryTrackParams
func NewQueryTrackParams(trackAddr string) QueryTrackParams {
	return QueryTrackParams{
		TrackAddr: trackAddr,
	}
}
