package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Sdp struct {
	From      sdk.AccAddress `json:"from" yaml:"from"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Data      []byte         `json:"data" yaml:"data"` // base64 data
}

func NewSdp(from sdk.AccAddress, recipient sdk.AccAddress, data []byte) *Sdp {
	return &Sdp{
		From:      from,
		Recipient: recipient,
		Data:      data,
	}
}

func (s Sdp) Validate() error {
	if s.From.Empty() {
		return fmt.Errorf(`from address cannot be empty`)
	}

	if s.Recipient.Empty() {
		return fmt.Errorf(`from address cannot be empty`)
	}

	if len(s.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	return nil
}
