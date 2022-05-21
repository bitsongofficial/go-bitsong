package types

import (
	"encoding/hex"
	"github.com/bitsongofficial/go-bitsong/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidProof(t *testing.T) {
	params.SetAddressPrefixes()

	address, err := sdk.AccAddressFromBech32("bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj2")
	if err != nil {
		assert.NoError(t, err)
	}

	amount := sdk.NewInt(1000000)

	root, _ := hex.DecodeString("5eb39dbca442a25db0f5d9e63489451b7bfc173796aa221e7207839de3a59e79")
	proofs := []string{
		"7f0b92cc8318e4fb0db9052325b474e2eabb80d79e6e1abab92093d3a88fe029",
		"a258c32bee9b0bbb7a2d1999ab4698294844e7440aa6dcd067e0d5142fa20522",
	}

	result := IsValidProof(uint64(0), address, amount, root, ConvertProofs(proofs))
	assert.True(t, result)
}
