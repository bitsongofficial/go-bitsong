package song

import (
	"github.com/BitSongOfficial/go-bitsong/x/song/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgSetTitle = types.NewMsgSetTitle
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgSetTitle      = types.MsgSetTitle
	QueryResTitles   = types.QueryResTitles
)