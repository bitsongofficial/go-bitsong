package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

func ConvertProofs(proofs []string) [][]byte {
	var proofsBz = make([][]byte, len(proofs))
	for i := 0; i < len(proofs); i++ {
		pBz, _ := hex.DecodeString(proofs[i])
		proofsBz[i] = pBz
	}

	return proofsBz
}

func IsValidProof(index uint64, account sdk.AccAddress, amount sdk.Int, root []byte, proofs [][]byte) bool {
	hasher := sha256.New()

	indexStr := strconv.FormatUint(index, 10)
	hashBz := crypto.Sha256([]byte(fmt.Sprintf("%s%s%s", indexStr, account.String(), amount.String())))

	for _, p := range proofs {
		hasher.Reset()
		if bytes.Compare(hashBz, p) < 0 {
			hasher.Write(hashBz)
			hasher.Write(p)
		} else {
			hasher.Write(p)
			hasher.Write(hashBz)
		}

		h := hasher.Sum(nil)
		hashBz = h
	}

	return bytes.Equal(hashBz, root)
}
