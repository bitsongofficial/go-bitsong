package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	MaxNameLength int = 140
)

type Distributor struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
}

func NewDistributor(name string, owner sdk.AccAddress) Distributor {
	return Distributor{Name: name, Address: owner}
}

func (d Distributor) String() string {
	return fmt.Sprintf(`Distributor: %s
  Address: %s`, d.Name, d.Address.String())
}

type Distributors []Distributor

func (d Distributors) String() string {
	out := "Name\n"

	for _, distr := range d {
		out += fmt.Sprintf("%s\n", distr.Name)
	}

	return strings.TrimSpace(out)
}
