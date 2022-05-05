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
	cdc.RegisterInterface((*FanTokenI)(nil), nil)

	cdc.RegisterConcrete(&FanToken{}, "go-bitsong/fantoken/FanToken", nil)

	cdc.RegisterConcrete(&MsgIssueFanToken{}, "go-bitsong/fantoken/MsgIssueFanToken", nil)
	cdc.RegisterConcrete(&MsgEditFanToken{}, "go-bitsong/fantoken/MsgEditFanToken", nil)
	cdc.RegisterConcrete(&MsgMintFanToken{}, "go-bitsong/fantoken/MsgMintFanToken", nil)
	cdc.RegisterConcrete(&MsgBurnFanToken{}, "go-bitsong/fantoken/MsgBurnFanToken", nil)
	cdc.RegisterConcrete(&MsgTransferFanTokenOwner{}, "go-bitsong/fantoken/MsgTransferFanTokenOwner", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssueFanToken{},
		&MsgEditFanToken{},
		&MsgMintFanToken{},
		&MsgBurnFanToken{},
		&MsgTransferFanTokenOwner{},
	)
	registry.RegisterInterface(
		"go-bitsong.fantoken.FanTokenI",
		(*FanTokenI)(nil),
		&FanToken{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
