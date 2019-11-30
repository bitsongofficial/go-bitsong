package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateArtist{}, "go-bitsong/MsgCreateArtist", nil)
	cdc.RegisterConcrete(MsgSetArtistImage{}, "go-bitsong/MsgSetArtistImage", nil)
	cdc.RegisterConcrete(MsgSetArtistStatus{}, "go-bitsong/MsgSetArtistStatus", nil)

	cdc.RegisterConcrete(ArtistVerifyProposal{}, "go-bitsong/ArtistVerifyProposal", nil)
}
