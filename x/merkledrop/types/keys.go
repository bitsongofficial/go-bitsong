package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "merkledrop"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

// Keys for merkledrop store
// Items are stored with the following keys => values
// - 0x01:<merkledropID_bytes>: merkledrop
// - 0x02:<owner>:<merkledropID_bytes>: merkledrop
// - 0x03: lastMerkledropID
// - 0x04:<merkledropID_bytes>:<merkledropIndex>: true
var (
	PrefixMerkleDrop        = []byte{0x01}
	PrefixMerkleDropByOwner = []byte{0x02}
	KeyLastMerkleDropId     = []byte{0x03}

	PrefixClaimedMerkleDrop = []byte{0x04}

	sep = []byte(":")
)

func MerkledropKey(id uint64) []byte {
	idBz := sdk.Uint64ToBigEndian(id)

	return genKey(PrefixMerkleDrop, sep, idBz)
}

func MerkledropOwnerKey(id uint64, owner sdk.AccAddress) []byte {
	idBz := sdk.Uint64ToBigEndian(id)
	return genKey(PrefixMerkleDropByOwner, sep, owner, sep, idBz)
}

func LastMerkledropIDKey() []byte {
	return KeyLastMerkleDropId
}

func ClaimedMerkledropKey(id, index uint64) []byte {
	return genKey(PrefixClaimedMerkleDrop, sep, sdk.Uint64ToBigEndian(id), sdk.Uint64ToBigEndian(index))
}

func genKey(bytes ...[]byte) (r []byte) {
	for _, b := range bytes {
		r = append(r, b...)
	}
	return
}
