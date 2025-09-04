package types

import (
	"bytes"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "nft"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	MaxDenomLength = 43
)

var (
	CollectionsPrefix      = collections.NewPrefix(0)
	SupplyPrefix           = collections.NewPrefix(1)
	NFTsPrefix             = collections.NewPrefix(2)
	NFTsByCollectionPrefix = collections.NewPrefix(3)
	NFTsByOwnerPrefix      = collections.NewPrefix(4)
)

func SplitNftLengthPrefixedKey(key []byte) (denom, tokenId []byte, err error) {
	parts := bytes.SplitN(key, []byte{0}, 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid composite key format: expected 2 parts, got %d", len(parts))
	}

	denomLen := len(parts[0])

	if denomLen > MaxDenomLength {
		return nil, nil, errors.Wrapf(sdkerrors.ErrInvalidType, "decoded denom key length %d exceeds max allowed length %d", denomLen, MaxDenomLength)
	}

	if len(key)-1 < denomLen {
		return nil, nil, fmt.Errorf("key is malformed: length prefix %d is greater than tokenId bytes %d", denomLen, len(key)-1)
	}

	denom = parts[0]
	tokenId = parts[1]

	return denom, tokenId, nil
}

func MustSplitNftLengthPrefixedKey(key []byte) (denom, tokenId []byte) {
	denom, tokenId, err := SplitNftLengthPrefixedKey(key)
	if err != nil {
		panic(err)
	}

	return denom, tokenId
}
