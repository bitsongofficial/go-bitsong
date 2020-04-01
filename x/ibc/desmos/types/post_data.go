package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
)

// Utility method that, given a songID, a creationTime and a sender address
// returns a new PostCreationData object
func NewSongCreationData(songID string, creationTime time.Time, sender sdk.AccAddress) posts.PostCreationData {
	return posts.NewPostCreationData(
		"",
		posts.PostID(0),
		true,
		DesmosBitsongSubspace,
		map[string]string{
			DesmosSongIDAttribute: songID,
		},
		sender,
		creationTime,
		nil,
		nil,
	)
}
