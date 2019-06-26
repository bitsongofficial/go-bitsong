package song

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Song tags
var (
	TxCategory = "song"

	Action         = sdk.TagAction
	Category       = sdk.TagCategory
	Owner          = sdk.TagSender
	SongId         = "id"
	//OrderStatus    = "status"
)
