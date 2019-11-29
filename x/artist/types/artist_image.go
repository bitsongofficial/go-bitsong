package types

type ArtistImage struct {
	Height uint64 `json:"height"`
	Width  uint64 `json:"width"`
	CID    string `json:"cid"`
}

// NewArtistImage returns an empty ArtistImage
func NewArtistImage(height uint64, width uint64, cid string) ArtistImage {
	return ArtistImage{
		Height: height,
		Width:  width,
		CID:    cid,
	}
}
