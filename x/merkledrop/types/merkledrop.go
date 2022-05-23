package types

import (
	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var (
	_ proto.Message = &Merkledrop{}
)

type MerkledropI interface {
	GetMerkleRoot() string
	GetAmount() string
}

/*func NewMerkledrop(merkleRoot string, amount sdk.Int, owner sdk.AccAddress) *Merkledrop {
	return &Merkledrop{
		MerkleRoot:  merkleRoot,
		TotalAmount: amount,
		Owner:       owner.String(),
	}
}*/

func (m Merkledrop) GetMerkleRoot() string {
	return m.MerkleRoot
}

func (m Merkledrop) GetAmount() string {
	return m.Amount.String()
}

func (m Merkledrop) String() string {
	bz, _ := yaml.Marshal(m)
	return string(bz)
}
