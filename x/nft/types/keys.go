package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "nft"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

var (
	PrefixNFT               = []byte{0x01}
	PrefixNFTByOwner        = []byte{0x02}
	PrefixMetadata          = []byte{0x03}
	PrefixCollection        = []byte{0x04}
	KeyPrefixLastMetadataId = []byte{0x05}
	KeyLastCollectionId     = []byte{0x06}
)

func CollectionKey(id uint64) []byte {
	return append(PrefixCollection, sdk.Uint64ToBigEndian(id)...)
}

func NftKey(idBz []byte) []byte {
	return append(PrefixNFT, idBz...)
}

func NftByOwnerKey(owner sdk.AccAddress, idBz []byte) []byte {
	return append(append(PrefixNFTByOwner, owner...), idBz...)
}

func LastMetadataId(collId uint64) []byte {
	return append(KeyPrefixLastMetadataId, sdk.Uint64ToBigEndian(collId)...)
}

func MetadataId(collId, metadataId uint64) []byte {
	return append(append(PrefixMetadata, sdk.Uint64ToBigEndian(collId)...), sdk.Uint64ToBigEndian(metadataId)...)
}
