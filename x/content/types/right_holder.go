package types

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type RightHolder struct {
	Quota   sdk.Dec        `json:"quota" yaml:"quota"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Rewards sdk.DecCoin    `json:"rewards" yaml:"rewards"`
}

func NewRightHolder(quota sdk.Dec, addr sdk.AccAddress) RightHolder {
	return RightHolder{
		Quota:   quota,
		Address: addr,
		Rewards: sdk.NewDecCoin(btsg.BondDenom, sdk.ZeroInt()),
	}
}

func (rh RightHolder) String() string {
	return fmt.Sprintf(`Quota: %s Address: %s Rewards:%s`, rh.Quota, rh.Address, rh.Rewards.String())
}

func (rh RightHolder) Equals(rh2 RightHolder) bool {
	return rh.Quota == rh2.Quota && rh.Address.Equals(rh2.Address) && rh.Rewards.IsEqual(rh2.Rewards)
}

func (rh RightHolder) Validate() error {
	if rh.Quota.IsNegative() {
		return fmt.Errorf(`right holder quota cannot be negative`)
	}

	if rh.Address.Empty() {
		return fmt.Errorf(`right holder address cannot be empty`)
	}

	if !rh.Rewards.IsEqual(sdk.DecCoin{}) {
		return fmt.Errorf(`right holder rewards must be empty`)
	}

	return nil
}

type RightsHolders []RightHolder

func (rhs RightsHolders) Validate() error {
	if len(rhs) > MaxRightsHoldersLength {
		return fmt.Errorf(`rights holders length cannot be greather then %d`, MaxRightsHoldersLength)
	}

	total := sdk.ZeroDec()
	for _, rh := range rhs {
		total = total.Add(rh.Quota)
	}

	if !total.Equal(sdk.NewDec(100)) {
		return fmt.Errorf(`sum of rights holders quota must be equal to 100`)
	}

	return nil
}

func (rhs RightsHolders) String() string {
	pnt := "Rights Holders:\nQuota Address\n"
	for _, rh := range rhs {
		pnt += fmt.Sprintf("%s %s\n", rh.Quota, rh.Address)
	}
	return strings.TrimSpace(pnt)
}

func (rhs RightsHolders) Equals(rhsCmp RightsHolders) bool {
	if len(rhs) != len(rhsCmp) {
		return false
	}

	for i, rh := range rhs {
		if !rh.Equals(rhsCmp[i]) {
			return false
		}
	}

	return true
}
