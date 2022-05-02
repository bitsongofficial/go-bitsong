package simulation

import (
	"bytes"
	"fmt"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixMetadata):
			var metaA, metaB types.Metadata
			cdc.MustUnmarshal(kvA.Value, &metaA)
			cdc.MustUnmarshal(kvB.Value, &metaB)
			return fmt.Sprintf("%v\n%v", metaA, metaB)
		case bytes.Equal(kvA.Key[:1], types.PrefixNFT):
			var nftA, nftB types.NFT
			cdc.MustUnmarshal(kvA.Value, &nftA)
			cdc.MustUnmarshal(kvB.Value, &nftB)
			return fmt.Sprintf("%v\n%v", nftA, nftB)
		case bytes.Equal(kvA.Key[:1], types.PrefixNFTByOwner):
			var nftA, nftB types.NFT
			cdc.MustUnmarshal(kvA.Value, &nftA)
			cdc.MustUnmarshal(kvB.Value, &nftB)
			return fmt.Sprintf("%v\n%v", nftA, nftB)
		case bytes.Equal(kvA.Key[:1], types.PrefixCollection):
			var collectionA, collectionB types.Collection
			cdc.MustUnmarshal(kvA.Value, &collectionA)
			cdc.MustUnmarshal(kvB.Value, &collectionB)
			return fmt.Sprintf("%v\n%v", collectionA, collectionB)
		case bytes.Equal(kvA.Key[:1], types.PrefixCollectionRecord):
			var recordA, recordB types.Collection
			cdc.MustUnmarshal(kvA.Value, &recordA)
			cdc.MustUnmarshal(kvB.Value, &recordB)
			return fmt.Sprintf("%v\n%v", recordA, recordB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
