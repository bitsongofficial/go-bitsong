package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Meta)(nil), nil)

	cdc.RegisterConcrete(MsgCreateArtist{}, "go-bitsong/MsgCreateArtist", nil)

	cdc.RegisterConcrete(GeneralMeta{}, "go-bitsong/ArtistGeneralMeta", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
