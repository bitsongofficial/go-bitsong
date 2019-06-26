package types

import (
	"fmt"
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Song is a struct that contains all the metadata of a title
type Song struct {
	Owner sdk.AccAddress `json:"owner"`
	Title string         `json:"title"`
}

func (s Song) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SONG:
  Owner:      %s
  Title: %s`,
		s.Owner,
		s.Title,
	))
}

type Songs []Song

func (songs Songs) String() string {
	out := ""
	for _, song := range songs {
		out += song.String() + "\n"
	}
	return out
}