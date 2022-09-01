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
	cdc.RegisterConcrete(&MsgCreateLaunchPad{}, "go-bitsong/launchpad/MsgCreateLaunchPad", nil)
	cdc.RegisterConcrete(&MsgUpdateLaunchPad{}, "go-bitsong/launchpad/MsgUpdateLaunchPad", nil)
	cdc.RegisterConcrete(&MsgCloseLaunchPad{}, "go-bitsong/launchpad/MsgCloseLaunchPad", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "go-bitsong/launchpad/MsgMintNFT", nil)
	cdc.RegisterConcrete(&MsgMintNFTs{}, "go-bitsong/launchpad/MsgMintNFTs", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateLaunchPad{},
		&MsgUpdateLaunchPad{},
		&MsgCloseLaunchPad{},
		&MsgMintNFT{},
		&MsgMintNFTs{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
