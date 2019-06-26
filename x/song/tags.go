package song

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Song tags
var (
	TxCategory = "song"

	Action                  = sdk.TagAction
	Category                = sdk.TagCategory
	Owner                   = sdk.TagSender
	SongID                  = "id"
	Content                 = "content"
	TotalReward             = "total_reward"
	RedistributionSplitRate = "redistribution_split_rate"
)
