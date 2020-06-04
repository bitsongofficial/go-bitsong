package types

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type DaoEntity struct {
	Shares  sdk.Dec        `json:"shares" yaml:"shares"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Rewards sdk.DecCoin    `json:"rewards" yaml:"rewards"`
}

func NewDaoEntity(shares sdk.Dec, addr sdk.AccAddress) DaoEntity {
	return DaoEntity{
		Shares:  shares,
		Address: addr,
		Rewards: sdk.NewDecCoin(btsg.BondDenom, sdk.ZeroInt()),
	}
}

func (de DaoEntity) String() string {
	return fmt.Sprintf(`Shares: %s Address: %s Rewards:%s`, de.Shares, de.Address, de.Rewards.String())
}

func (de DaoEntity) Equals(de2 DaoEntity) bool {
	return de.Shares == de2.Shares && de.Address.Equals(de2.Address) && de.Rewards.IsEqual(de2.Rewards)
}

func (de DaoEntity) Validate() error {
	if de.Shares.IsNegative() {
		return fmt.Errorf(`dao entity shares cannot be negative`)
	}

	if de.Address.Empty() {
		return fmt.Errorf(`dao entity address cannot be empty`)
	}

	if !de.Rewards.IsEqual(sdk.DecCoin{}) {
		return fmt.Errorf(`dao entity rewards must be empty`)
	}

	return nil
}

type Dao []DaoEntity

func (d Dao) Validate() error {
	if len(d) > MaxDaoLength {
		return fmt.Errorf(`dao length cannot be greather then %d`, MaxDaoLength)
	}

	return nil
}

func (d Dao) String() string {
	pnt := "Dao Entity:\nShares Address\n"
	for _, de := range d {
		pnt += fmt.Sprintf("%s %s\n", de.Shares, de.Address)
	}
	return strings.TrimSpace(pnt)
}

func (d Dao) Equals(dCmp Dao) bool {
	if len(d) != len(dCmp) {
		return false
	}

	for i, de := range d {
		if !de.Equals(dCmp[i]) {
			return false
		}
	}

	return true
}
