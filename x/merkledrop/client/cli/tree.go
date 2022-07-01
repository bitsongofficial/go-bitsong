package cli

import (
	"bytes"
	"crypto/sha256"
)

type Tree [][][]byte

// NewTree builds and returns a new merkle tree with the given leaf nodes.
// Each input leaf is hashed before the tree is built
func NewTree(leafs ...[]byte) Tree {
	hasher := sha256.New()

	hashed := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		hasher.Reset()
		hasher.Write(leaf)
		hashed[i] = hasher.Sum(nil)
	}

	tree := make([][][]byte, 1)
	tree[0] = hashed

	for {

		row := tree[len(tree)-1]
		l := len(row)

		if l == 1 {
			break
		}

		if l%2 != 0 {
			l--
		}

		var next [][]byte

		for i := 0; i < l; i += 2 {
			hasher.Reset()
			// Smaller one first, otherwise proofs will not work
			if bytes.Compare(row[i], row[i+1]) < 0 {
				hasher.Write(row[i])
				hasher.Write(row[i+1])
			} else {
				hasher.Write(row[i+1])
				hasher.Write(row[i])
			}
			h := hasher.Sum(nil)
			next = append(next, h)
		}

		// Append to current row on odd count
		if l != len(row) {
			next = append(next, row[l])
		}

		tree = append(tree, next)
	}

	return tree
}

// Root returns the root of the merkle tree i.e the merkle root
func (t Tree) Root() []byte {
	return t[len(t)-1][0]
}

func (t Tree) Height() int {
	return len(t)
}

// Leafs returns the leafs in the tree. This is the hashes of the original
// input data
func (t Tree) Leafs() [][]byte {
	return t[0]
}

// LeafIndex returns the index of the given leaf and a -1 if the leaf is not
// found
func (t Tree) LeafIndex(leaf []byte) int {
	for i, item := range t[0] {
		if bytes.Equal(item, leaf) {
			return i
		}
	}
	return -1
}

// Proof returns the proof of a given leaf node.
func (t Tree) Proof(leaf []byte) (out [][]byte) {
	idx := t.LeafIndex(leaf)
	if idx < 0 {
		return
	}

	// Exclude row with root hash i.e last row
	li := len(t) - 1
	for _, row := range t[:li] {
		if (idx)%2 == 0 {
			if idx < len(row)-1 {
				out = append(out, row[idx+1])
			}
			// If this is the last odd element we do nothing as it is carried up.
		} else {
			out = append(out, row[idx-1])
		}
		// Cut in half for the evaluation of the next row
		idx = idx / 2
	}

	out = append(out, t[li][0])

	return
}
