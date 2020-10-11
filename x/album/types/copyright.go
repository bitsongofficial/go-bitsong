package types

type Copyright struct {
	Text string `json:"text" yaml:"text"` // The copyright text for the album.
	Type string `json:"type" yaml:"type"` // The type of copyright.
}
