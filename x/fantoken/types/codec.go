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
	cdc.RegisterConcrete(&MsgIssue{}, "go-bitsong/fantoken/MsgIssue", nil)
	cdc.RegisterConcrete(&MsgEdit{}, "go-bitsong/fantoken/MsgEdit", nil)
	cdc.RegisterConcrete(&MsgMint{}, "go-bitsong/fantoken/MsgMint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "go-bitsong/fantoken/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgTransferOwnership{}, "go-bitsong/fantoken/MsgTransferOwnership", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssue{},
		&MsgEdit{},
		&MsgMint{},
		&MsgBurn{},
		&MsgTransferOwnership{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
