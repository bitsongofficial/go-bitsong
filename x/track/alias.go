package track

import (
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
	DefaultStartingSongID = types.DefaultStartingSongID
)

var (
	NewMsgPublish = types.NewMsgPublish
	NewMsgPlay    = types.NewMsgPlay
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgPublish      = types.MsgPublish
	MsgPlay         = types.MsgPlay
	Song            = types.Song
	Songs			= types.Songs
	Pool			= types.Pool
)