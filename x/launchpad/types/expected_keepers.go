package types

import (
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error

	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

type NftKeeper interface {
	GetCollectionById(ctx sdk.Context, id uint64) (nfttypes.Collection, error)
	SetCollection(ctx sdk.Context, collection nfttypes.Collection)
	GetNFTById(ctx sdk.Context, id string) (nfttypes.NFT, error)
	GetMetadataById(ctx sdk.Context, collId, id uint64) (nfttypes.Metadata, error)
	TransferNFT(ctx sdk.Context, msg *nfttypes.MsgTransferNFT) error
	UpdateMetadataAuthority(ctx sdk.Context, msg *nfttypes.MsgUpdateMetadataAuthority) error
	UpdateMintAuthority(ctx sdk.Context, msg *nfttypes.MsgUpdateMintAuthority) error
	SetPrimarySaleHappened(ctx sdk.Context, collId, metadataId uint64) error
	PrintEdition(ctx sdk.Context, msg *nfttypes.MsgPrintEdition) (string, error)
	CreateNFT(ctx sdk.Context, msg *nfttypes.MsgCreateNFT) (uint64, string, error)
	SetMetadata(ctx sdk.Context, metadata nfttypes.Metadata)
	SetNFT(ctx sdk.Context, nft nfttypes.NFT)
	GetLastMetadataId(ctx sdk.Context, collId uint64) uint64
	SetLastMetadataId(ctx sdk.Context, collId, id uint64)
}
