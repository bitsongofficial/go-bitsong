package types

// Constants pertaining to a Meta object
const (
	MaxNameLength int = 140
)

type Meta struct {
	Name string `json:"name" yaml:"name"`
}

func NewMeta(name string) Meta {
	//return Meta{name, images}
	return Meta{name}
}

// Images
type Image struct {
	CID    string
	Height string
	Width  string
}

func NewImage(cid string, width string, height string) Image {
	return Image{
		CID:    cid,
		Height: height,
		Width:  width,
	}
}
