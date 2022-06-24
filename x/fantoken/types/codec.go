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
	cdc.RegisterConcrete(&MsgIssue{}, "go-bitsong/fantoken/MsgIssue", nil)
	cdc.RegisterConcrete(&MsgMint{}, "go-bitsong/fantoken/MsgMint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "go-bitsong/fantoken/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgDisableMint{}, "go-bitsong/fantoken/MsgDisableMint", nil)
	cdc.RegisterConcrete(&MsgSetAuthority{}, "go-bitsong/fantoken/MsgSetAuthority", nil)
	cdc.RegisterConcrete(&MsgSetMinter{}, "go-bitsong/fantoken/MsgSetMinter", nil)
	cdc.RegisterConcrete(&MsgSetUri{}, "go-bitsong/fantoken/MsgSetUri", nil)
	cdc.RegisterConcrete(&UpdateFeesProposal{}, "go-bitsong/fantoken/UpdateFeesProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssue{},
		&MsgMint{},
		&MsgBurn{},
		&MsgDisableMint{},
		&MsgSetAuthority{},
		&MsgSetMinter{},
		&MsgSetUri{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateFeesProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
