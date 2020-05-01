package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddContent{}, "go-bitsong/MsgAddContent", nil)
	cdc.RegisterConcrete(MsgMintContent{}, "go-bitsong/MsgMintContent", nil)
	cdc.RegisterConcrete(MsgBurnContent{}, "go-bitsong/MsgBurnContent", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
