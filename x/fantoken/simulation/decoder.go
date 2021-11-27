package simulation

// DONTCOVER

import (
	"bytes"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/bitsongofficial/chainmodules/x/fantoken/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding token type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PrefixFanTokenForDenom):
			var tokenA, tokenB types.FanToken
			cdc.MustUnmarshalBinaryBare(kvA.Value, &tokenA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &tokenB)
			return fmt.Sprintf("%v\n%v", tokenA, tokenB)
		case bytes.Equal(kvA.Key[:1], types.PrefixFanTokens):
			var denomA, denomB gogotypes.Value
			cdc.MustUnmarshalBinaryBare(kvA.Value, &denomA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &denomB)
			return fmt.Sprintf("%v\n%v", denomA, denomB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
