package types

import (
	"github.com/bitsongofficial/go-bitsong/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidProof(t *testing.T) {
	params.SetAddressPrefixes()

	address, err := sdk.AccAddressFromBech32("bitsong1qyc9ccfx7clj0kswgjz6cdf5f4n6v4nycc3dha")
	if err != nil {
		assert.NoError(t, err)
	}

	amount := sdk.NewInt(1000000)

	proofs := []string{
		"20245fe3fcdbf17069bc0de04e319296766a7138be5e5a27c6f5bc05e0c23de9",
		"b8fedba5a18186d4fb92ffcf9924b408d6048aaeb76b10cad97cf6be4071b710",
	}

	result := IsValidProof(address, amount, ConvertProofs(proofs))
	assert.True(t, result)
}
