package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
)

const (
	DesmosBitsongSubspace = "a31be8a1946fb15200d7081163bf3c41eae3b8b745e8bbf7d96e04e57c9ddf9b"
	DesmosSongIDAttribute = "song_id"
)

// Utility method that, given a songID, a creationTime and a sender address
// returns a new PostCreationData object
func NewSongCreationData(songID string, creationTime time.Time, sender sdk.AccAddress) posts.PostCreationData {
	return posts.NewPostCreationData(
		songID,
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
