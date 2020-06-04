package types

type Image struct {
	Height uint   `json:"height" yaml:"height"`
	Width  uint   `json:"width" yaml:"width"`
	Url    string `json:"url" yaml:"url"`
}
