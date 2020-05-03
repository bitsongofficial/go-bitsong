package types

const (
	QueryParams = "params"
	QueryID     = "id"
)

// Params for queries
type QueryPlayerParams struct {
	ID uint64 `json:"id" yaml:"id"`
}

// creates a new instance of QueryPlayerParams
func NewQueryPlayerParams(id uint64) QueryPlayerParams {
	return QueryPlayerParams{
		ID: id,
	}
}
