package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the IBC transfer codec.
var ModuleCdc = codec.New()

// RegisterCodec registers the IBC transfer types
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreteSongPost{}, "bitsong/ibc/MsgCreateSongPost", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
