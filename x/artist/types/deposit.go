package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deposit
type Deposit struct {
	ArtistID  uint64         `json:"artist_id" yaml:"artist_id"` // the artist id
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"` //  Address of the depositor
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`       //  Deposit amount
}

// NewDeposit creates a new Deposit instance
func NewDeposit(artistID uint64, depositor sdk.AccAddress, amount sdk.Coins) Deposit {
	return Deposit{
		artistID,
		depositor,
		amount,
	}
}

func (d Deposit) String() string {
	return fmt.Sprintf("deposit by %s on ArtistID %d is for the amount %s",
		d.Depositor, d.ArtistID, d.Amount)
}

// Deposits is a collection of Deposit objects
type Deposits []Deposit

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for ArtistID %d:", d[0].ArtistID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, dep.Amount)
	}
	return out
}

// Equals returns whether two deposits are equal.
func (d Deposit) Equals(comp Deposit) bool {
	return d.Depositor.Equals(comp.Depositor) && d.ArtistID == comp.ArtistID && d.Amount.IsEqual(comp.Amount)
}

// Empty returns whether a deposit is empty.
func (d Deposit) Empty() bool {
	return d.Equals(Deposit{})
}
