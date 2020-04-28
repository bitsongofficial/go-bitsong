package types

import (
	"fmt"

	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

const (
	AttributeKeySongID = "song_id"
)

// IBC transfer events vars
var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, ModuleName)
)
