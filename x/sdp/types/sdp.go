package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Sdp struct {
	Hash      string         `json:"hash" yaml:"hash"`
	From      sdk.AccAddress `json:"from" yaml:"from"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Offer     []byte         `json:"data" yaml:"data"`    // base64 data
	Answer    []byte         `json:"answer" yaml:"anser"` // base64 data
	Expire    time.Time      `json:"expire" yaml:"expire"`
}

func NewSdp(hash string, from sdk.AccAddress, recipient sdk.AccAddress, offer []byte, answer []byte) *Sdp {
	return &Sdp{
		Hash:      hash,
		From:      from,
		Recipient: recipient,
		Offer:     offer,
		Answer:    answer,
		Expire:    time.Now().Add(5 * time.Minute),
	}
}

func (s Sdp) Validate() error {
	if s.From.Empty() {
		return fmt.Errorf(`from address cannot be empty`)
	}

	if s.Recipient.Empty() {
		return fmt.Errorf(`from address cannot be empty`)
	}

	if len(s.Offer) == 0 && len(s.Answer) == 0 {
		return fmt.Errorf("offer or answer cannot be empty")
	}

	return nil
}
