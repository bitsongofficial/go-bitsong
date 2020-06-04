package types

const (
	QueryParams = "params"
	QueryCid    = "cid"
)

// Params for queries
type QueryTrackParams struct {
	Cid string `json:"cid" yaml:"cid"`
}

// creates a new instance of QueryContentParams
func NewQueryContentParams(cid string) QueryTrackParams {
	return QueryTrackParams{
		Cid: cid,
	}
}
