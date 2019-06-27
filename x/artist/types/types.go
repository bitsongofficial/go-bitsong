package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Artist is a struct that contains all the metadata of a title
type Artist struct {
	ArtistID string         `json:"artist_id"`
	Image    string         `json:"image"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
}

func (a Artist) String() string {
	return fmt.Sprintf(`Song %d:
		  Artist ID:	%s
		  Image:		%s
		  Name:			%s`, a.ArtistID, a.Image, a.Name)
}

// Artists is an array of song
// To FIX with new fields
type Artists []*Artist

func (songs Artists) String() string {
	out := fmt.Sprintf("%10s - (%15s) - (%40s) - [%10s] - Create Time\n", "ID", "Title", "Owner", "CreateTime")
	for _, song := range songs {
		out += fmt.Sprintf("%10d - (%15s) - (%40s) - [%10s]\n",
			a.ArtistID, a.Image, a.Name)
	}

	return strings.TrimSpace(out)
}
