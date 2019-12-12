package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName        = "artist"   // ModuleName is the name of the module
	StoreKey          = ModuleName // StoreKey is the store key string for artist
	RouterKey         = ModuleName // RouterKey is the message route for artist
	DefaultParamspace = ModuleName // DefaultParamspace default name for parameter store
)

// Keys for artist store
// Items are stored with the following key: values
//
// - 0x00<artistID_Bytes>: Artist
//
// - 0x01: nextArtistID
//
// - 0x10<artistID_Bytes><depositorAddr_Bytes>: Deposit
var (
	ArtistsKeyPrefix = []byte{0x00}
	ArtistIDKey      = []byte{0x01}

	DepositsKeyPrefix = []byte{0x10}
)

// DepositsKey gets the first part of the deposits key based on the artistID
func DepositsKey(artistID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, artistID)
	return append(DepositsKeyPrefix, bz...)
}

// DepositKey key of a specific deposit from the store
func DepositKey(artistID uint64, depositorAddr sdk.AccAddress) []byte {
	return append(DepositsKey(artistID), depositorAddr.Bytes()...)
}

// SplitKeyDeposit split the deposits key and returns the artistID and depositor address
/*func SplitKeyDeposit(key []byte) (artistID uint64, depositorAddr sdk.AccAddress) {
	return splitKeyWithAddress(key)
}

func splitKeyWithAddress(key []byte) (artistID uint64, addr sdk.AccAddress) {
	if len(key[1:]) != 8+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8+sdk.AddrLen))
	}

	artistID = binary.LittleEndian.Uint64(key[1:9])
	addr = sdk.AccAddress(key[9:])
	return
}*/
