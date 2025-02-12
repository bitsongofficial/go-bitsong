package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino    = codec.NewLegacyAmino()
	AminoCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

}

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterCadanceContract{}, "cadance/MsgRegisterCadanceContract", nil)
	cdc.RegisterConcrete(&MsgUnregisterCadanceContract{}, "cadance/MsgUnregisterCadanceContract", nil)
	cdc.RegisterConcrete(&MsgUnjailCadanceContract{}, "cadance/MsgUnjailCadanceContract", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "cadance/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterCadanceContract{},
		&MsgUnregisterCadanceContract{},
		&MsgUnjailCadanceContract{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
