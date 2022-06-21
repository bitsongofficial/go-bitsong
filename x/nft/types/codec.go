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
	cdc.RegisterConcrete(&MsgCreateNFT{}, "go-bitsong/nft/MsgCreateNFT", nil)
	cdc.RegisterConcrete(&MsgTransferNFT{}, "go-bitsong/nft/MsgTransferNFT", nil)
	cdc.RegisterConcrete(&MsgSignMetadata{}, "go-bitsong/nft/MsgSignMetadata", nil)
	cdc.RegisterConcrete(&MsgUpdateMetadata{}, "go-bitsong/nft/MsgUpdateMetadata", nil)
	cdc.RegisterConcrete(&MsgUpdateMetadataAuthority{}, "go-bitsong/nft/MsgUpdateMetadataAuthority", nil)
	cdc.RegisterConcrete(&MsgCreateCollection{}, "go-bitsong/nft/MsgCreateCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateCollectionAuthority{}, "go-bitsong/nft/MsgUpdateCollectionAuthority", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateNFT{},
		&MsgTransferNFT{},
		&MsgSignMetadata{},
		&MsgUpdateMetadata{},
		&MsgUpdateMetadataAuthority{},
		&MsgCreateCollection{},
		&MsgUpdateCollectionAuthority{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
