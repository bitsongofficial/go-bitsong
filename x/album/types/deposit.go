package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deposit
type Deposit struct {
	AlbumID   uint64         `json:"album_id" yaml:"album_id"`   // the album id
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"` //  Address of the depositor
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`       //  Deposit amount
}

// NewDeposit creates a new Deposit instance
func NewDeposit(albumID uint64, depositor sdk.AccAddress, amount sdk.Coins) Deposit {
	return Deposit{
		albumID,
		depositor,
		amount,
	}
}

func (d Deposit) String() string {
	return fmt.Sprintf("deposit by %s on AlbumID %d is for the amount %s",
		d.Depositor, d.AlbumID, d.Amount)
}

// Deposits is a collection of Deposit objects
type Deposits []Deposit

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for AlbumID %d:", d[0].AlbumID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, dep.Amount)
	}
	return out
}

// Equals returns whether two deposits are equal.
func (d Deposit) Equals(comp Deposit) bool {
	return d.Depositor.Equals(comp.Depositor) && d.AlbumID == comp.AlbumID && d.Amount.IsEqual(comp.Amount)
}

// Empty returns whether a deposit is empty.
func (d Deposit) Empty() bool {
	return d.Equals(Deposit{})
}
