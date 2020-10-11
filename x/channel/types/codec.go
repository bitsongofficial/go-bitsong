package types

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgChannelCreate{}, "go-bitsong/MsgChannelCreate", nil)
	cdc.RegisterConcrete(MsgChannelEdit{}, "go-bitsong/MsgChannelEdit", nil)
	cdc.RegisterConcrete(Channel{}, "go-bitsong/Channel", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
