package types

import (
	"fmt"
)

func (nft NFT) Id() string {
	return fmt.Sprintf("%d:%d:%d", nft.CollId, nft.MetadataId, nft.Seq)
}
