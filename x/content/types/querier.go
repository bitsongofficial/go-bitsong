package types

const (
	QueryParams  = "params"
	QueryContent = "content"
)

// Params for queries
type QueryContentParams struct {
	Uri string `json:"uri" yaml:"uri"`
}

// creates a new instance of QueryContentParams
func NewQueryContentParams(uri string) QueryContentParams {
	return QueryContentParams{
		Uri: uri,
	}
}
