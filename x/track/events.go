package track

import (
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
)

// Track module event types
var (
	EventTypePublish = "publish"

	AttributeKeyTrackId                 = "track_id"
	AttributeKeyTitle                   = "title"
	AttributeKeyContent                 = "content"
	AttributeKeyRedistributionSplitRate = "redistribution_split_rate"

	AttributeValueCategory = types.ModuleName
)
