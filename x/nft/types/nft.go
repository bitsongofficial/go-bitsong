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

func NftIdToBytes(id string) ([]byte, error) {
	splits := strings.Split(id, ":")
	if len(splits) != 3 {
		return nil, fmt.Errorf("invalid nft id: %s", id)
	}
	collId, err := strconv.Atoi(splits[0])
	if err != nil || collId < 0 {
		return nil, fmt.Errorf("invalid nft id: %s", id)
	}
	metadataId, err := strconv.Atoi(splits[1])
	if err != nil || metadataId < 0 {
		return nil, fmt.Errorf("invalid nft id: %s", id)
	}
	seq, err := strconv.Atoi(splits[2])
	if err != nil || seq < 0 {
		return nil, fmt.Errorf("invalid nft id: %s", id)
	}

	return (NFT{
		CollId:     uint64(collId),
		MetadataId: uint64(metadataId),
		Seq:        uint64(seq),
	}).IdBytes(), nil
}
