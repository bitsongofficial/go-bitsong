package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProof2(t *testing.T) {
	leafs := [][]byte{
		[]byte("2bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw3000000"),
		[]byte("0bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj21000000"),
		[]byte("1bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw2000000"),
	}

	tree := NewTree(leafs...)
	merkleRootStr := fmt.Sprintf("%x", tree.Root())
	assert.Equal(t, "5eb39dbca442a25db0f5d9e63489451b7bfc173796aa221e7207839de3a59e79", merkleRootStr)
}

func TestCreateProof(t *testing.T) {
	params.SetAddressPrefixes()

	accounts := map[string]string{
		"bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj2": "1000000ubtsg",
		"bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw": "2000000ubtsg",
		"bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw": "3000000ubtsg",
	}

	accMap, err := AccountsFromMap(accounts)
	assert.NoError(t, err)

	tree, claimInfo, err := CreateDistributionList(accMap)
	assert.NoError(t, err)

	for _, l := range tree.Leafs() {
		fmt.Println(fmt.Sprintf("%x", l))
		fmt.Println(fmt.Sprintf("proof: %x", tree.Proof(l)))
	}

	merkleRoot := fmt.Sprintf("%x", tree.Root())
	assert.Equal(t, "5eb39dbca442a25db0f5d9e63489451b7bfc173796aa221e7207839de3a59e79", merkleRoot)

	fmt.Println(merkleRoot)

	for i, c := range claimInfo {
		fmt.Println(fmt.Sprintf("%s %s: ", i, c.Proof))
	}
}
