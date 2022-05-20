package keeper

import (
	"fmt"
	"testing"
)

var (
	claimedBitMap = make(map[uint64]uint64)
)

func isClaimed(index uint64) bool {
	claimedWordIndex := index / 256
	claimedBitIndex := index % 256
	fmt.Println(claimedBitMap[claimedWordIndex])
	claimedWord := claimedBitMap[claimedWordIndex]
	mask := uint64(1 << claimedBitIndex)

	fmt.Println(claimedWord & mask)
	fmt.Println(claimedWord)
	fmt.Println(mask)

	return claimedWord&mask == mask
}

func setClaim(index uint64) {
	claimedWordIndex := index / 256
	claimedBitIndex := index % 256
	claimedBitMap[claimedWordIndex] = claimedBitMap[claimedWordIndex] | (1 << claimedBitIndex)
}

func TestKeeper_SetClaimed(t *testing.T) {
	index := uint64(253265485458494684)

	isClaim := isClaimed(index)
	fmt.Println(fmt.Sprintf("is claim %v", isClaim))

	setClaim(index)

	fmt.Println(fmt.Sprintf("is claim %v", isClaimed(index)))
}
