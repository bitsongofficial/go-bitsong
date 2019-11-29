package types

import (
	"fmt"

	c "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
)

type ArtistImage struct {
	Height uint64 `json:"height"`
	Width  uint64 `json:"width"`
	CID    string `json:"cid"`
}

func (image ArtistImage) Validate() error {
	if fmt.Sprintf("%d", image.Height) == "" || image.Height < 0 {
		return fmt.Errorf("invalid artist image height: %d", image.Height)
	}

	if fmt.Sprintf("%d", image.Width) == "" || image.Width < 0 {
		return fmt.Errorf("invalid artist image width: %d", image.Width)
	}

	_, err := c.Decode(image.CID)
	if err != nil {
		return fmt.Errorf("invalid artist image cid: %s", image.CID)
	}

	return nil
}

// NewArtistImage returns an empty ArtistImage
func NewArtistImage(height uint64, width uint64, cid string) ArtistImage {
	return ArtistImage{
		Height: height,
		Width:  width,
		CID:    cid,
	}
}
