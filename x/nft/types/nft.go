package types

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (nft NFT) Id() string {
	return fmt.Sprintf("%d:%d:%d", nft.CollId, nft.MetadataId, nft.Seq)
}

func (nft NFT) IdBytes() []byte {
	return append(append(sdk.Uint64ToBigEndian(nft.CollId), sdk.Uint64ToBigEndian(nft.MetadataId)...), sdk.Uint64ToBigEndian(nft.Seq)...)
}

func IsValidNftId(id string) bool {
	splits := strings.Split(id, ":")
	collId, err := strconv.Atoi(splits[0])
	if err != nil || collId < 0 {
		return false
	}
	metadataId, err := strconv.Atoi(splits[1])
	if err != nil || metadataId < 0 {
		return false
	}
	seq, err := strconv.Atoi(splits[2])
	if err != nil || seq < 0 {
		return false
	}
	return true
}

func NftIdToBytes(id string) []byte {
	splits := strings.Split(id, ":")
	if !IsValidNftId(id) {
		panic("invalid nft id")
	}
	collId, _ := strconv.Atoi(splits[0])
	metadataId, _ := strconv.Atoi(splits[1])
	seq, _ := strconv.Atoi(splits[2])
	return (NFT{
		CollId:     uint64(collId),
		MetadataId: uint64(metadataId),
		Seq:        uint64(seq),
	}).IdBytes()
}
