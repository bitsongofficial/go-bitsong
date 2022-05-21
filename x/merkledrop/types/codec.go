package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*MerkledropI)(nil), nil)
	cdc.RegisterConcrete(&Merkledrop{}, "go-bitsong/merkledrop/Merkledrop", nil)

	cdc.RegisterConcrete(&MsgCreate{}, "go-bitsong/merkledrop/MsgCreate", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "go-bitsong/merkledrop/MsgClaim", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreate{},
		&MsgClaim{},
	)
	registry.RegisterInterface(
		"go-bitsong.merkledrop.MerkledropI",
		(*MerkledropI)(nil),
		&Merkledrop{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
