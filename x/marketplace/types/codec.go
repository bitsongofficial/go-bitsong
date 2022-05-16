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
	cdc.RegisterConcrete(&MsgCreateAuction{}, "go-bitsong/nft/MsgCreateAuction", nil)
	cdc.RegisterConcrete(&MsgSetAuctionAuthority{}, "go-bitsong/nft/MsgSetAuctionAuthority", nil)
	cdc.RegisterConcrete(&MsgStartAuction{}, "go-bitsong/nft/MsgStartAuction", nil)
	cdc.RegisterConcrete(&MsgEndAuction{}, "go-bitsong/nft/MsgEndAuction", nil)
	cdc.RegisterConcrete(&MsgPlaceBid{}, "go-bitsong/nft/MsgPlaceBid", nil)
	cdc.RegisterConcrete(&MsgCancelBid{}, "go-bitsong/nft/MsgCancelBid", nil)
	cdc.RegisterConcrete(&MsgClaimBid{}, "go-bitsong/nft/MsgClaimBid", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAuction{},
		&MsgSetAuctionAuthority{},
		&MsgStartAuction{},
		&MsgEndAuction{},
		&MsgPlaceBid{},
		&MsgCancelBid{},
		&MsgClaimBid{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
