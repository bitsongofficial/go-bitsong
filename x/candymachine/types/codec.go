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
	cdc.RegisterConcrete(&MsgCreateCandyMachine{}, "go-bitsong/candymachine/MsgCreateCandyMachine", nil)
	cdc.RegisterConcrete(&MsgUpdateCandyMachine{}, "go-bitsong/candymachine/MsgUpdateCandyMachine", nil)
	cdc.RegisterConcrete(&MsgCloseCandyMachine{}, "go-bitsong/candymachine/MsgCloseCandyMachine", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "go-bitsong/candymachine/MsgMintNFT", nil)
	cdc.RegisterConcrete(&MsgMintNFTs{}, "go-bitsong/candymachine/MsgMintNFTs", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCandyMachine{},
		&MsgUpdateCandyMachine{},
		&MsgCloseCandyMachine{},
		&MsgMintNFT{},
		&MsgMintNFTs{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
