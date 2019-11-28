package types

// Constants pertaining to a Meta object
const (
	MaxNameLength int = 140
)

type Meta struct {
	Name string `json:"name" yaml:"name"`
	//Images []Image `json:"images" yaml:"images"`
}

func NewMeta(name string) Meta {
	//return Meta{name, images}
	return Meta{name}
}

// Images
type Image struct {
	CID    string
	Height uint64
	Width  uint64
}

func NewImage(cid string, width uint64, height uint64) Image {
	return Image{
		CID:    cid,
		Height: height,
		Width:  width,
	}
}
