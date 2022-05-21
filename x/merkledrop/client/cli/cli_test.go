package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProof2(t *testing.T) {
	leafs := [][]byte{
		[]byte("0bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj21000000"),
		[]byte("1bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw2000000"),
		[]byte("2bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw3000000"),
	}

	tree := NewTree(leafs...)
	merkleRootStr := fmt.Sprintf("%x", tree.Root())
	assert.Equal(t, "3452cae72dab475d017c1c46d289f9dc458a9fccf79add3e49347f2fc984e463", merkleRootStr)

	/*fmt.Println(fmt.Sprintf("%x", tree.Proof(crypto.Sha256(leafs[0]))))
	fmt.Println(fmt.Sprintf("%x", tree.Proof(crypto.Sha256(leafs[1]))))
	fmt.Println(fmt.Sprintf("%x", tree.Proof(crypto.Sha256(leafs[2]))))*/
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

	tree, _, err := CreateDistributionList(accMap)
	assert.NoError(t, err)

	merkleRoot := fmt.Sprintf("%x", tree.Root())
	assert.Equal(t, "3452cae72dab475d017c1c46d289f9dc458a9fccf79add3e49347f2fc984e463", merkleRoot)
}
