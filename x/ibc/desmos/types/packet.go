package types

import (
	"time"

	ibcposts "github.com/desmos-labs/desmos/x/ibc/posts"
)

const (
	DesmosBitsongSubspace = "a31be8a1946fb15200d7081163bf3c41eae3b8b745e8bbf7d96e04e57c9ddf9b"
	DesmosSongIDAttribute = "song_id"
)

// Utility method that, given a songID, a creationTime and a sender address
// returns a new PostCreationData object
func NewSongCreationData(songID string, creationTime time.Time, postOwner string) ibcposts.PostCreationPacketData {
	return ibcposts.NewPostCreationPacketData(
		songID,
		0,
		true,
		DesmosBitsongSubspace,
		map[string]string{
			DesmosSongIDAttribute: songID,
		},
		postOwner,
		creationTime,
		nil,
		nil,
	)
}
