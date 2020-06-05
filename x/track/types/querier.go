package types

const (
	QueryParams = "params"
	QueryID     = "id"
)

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
