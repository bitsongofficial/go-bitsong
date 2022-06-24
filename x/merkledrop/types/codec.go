package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	cdc.RegisterConcrete(&MsgCreate{}, "go-bitsong/merkledrop/MsgCreate", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "go-bitsong/merkledrop/MsgClaim", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "go-bitsong/merkledrop/MsgWithdraw", nil)
	cdc.RegisterConcrete(&UpdateFeesProposal{}, "go-bitsong/merkledrop/UpdateFeesProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreate{},
		&MsgClaim{},
		&MsgWithdraw{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateFeesProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
