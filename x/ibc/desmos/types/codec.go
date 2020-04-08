package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	dsmtypes "github.com/desmos-labs/desmos/x/posts"
)

// ModuleCdc defines the IBC transfer codec.
var ModuleCdc = codec.New()

// RegisterCodec registers the IBC transfer types
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateSongPost{}, "bitsong/ibc/MsgCreateSongPost", nil)
	dsmtypes.RegisterCodec(cdc)
}

func init() {
	RegisterCodec(ModuleCdc)
}
