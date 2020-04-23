package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type RightHolder struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Quota   uint           `json:"quota" yaml:"quota"`
}

func NewRightHolder(addr sdk.AccAddress, quota uint) RightHolder {
	return RightHolder{
		Address: addr,
		Quota:   quota,
	}
}

type RightsHolders []RightHolder

func (rhs RightsHolders) String() string {
	out := "Quota - Address\n"
	for _, rh := range rhs {
		out += fmt.Sprintf("%d - %s\n", rh.Quota, rh.Address.String())
	}
	return strings.TrimSpace(out)
}

func (rhs RightsHolders) Equals(rightsHolders RightsHolders) bool {
	// TODO: need better logic
	return len(rhs) == len(rightsHolders)
}

func (rhs RightsHolders) Sum() uint {
	sum := uint(0)
	for _, rh := range rhs {
		sum += rh.Quota
	}
	return sum
}

func (rhs RightsHolders) Validate() error {
	if rhs.Sum() != 100 {
		return fmt.Errorf("rights holders quota must be equal to 100")
	}

	return nil
}
