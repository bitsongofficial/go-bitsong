package types

import "github.com/cosmos/cosmos-sdk/codec"

// ModuleCdc is the codec
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateTrack{}, "go-bitsong/MsgCreateTrack", nil)
	cdc.RegisterConcrete(MsgPlay{}, "go-bitsong/MsgPlay", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "go-bitsong/MsgDepositTrack", nil)
}
