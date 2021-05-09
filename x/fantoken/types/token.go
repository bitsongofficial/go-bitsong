package types

import (
	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ proto.Message = &FanToken{}
)

// TokenI defines an interface for Token
type FanTokenI interface {
	GetDenom() string
	GetName() string
	GetMaxSupply() sdk.Int
	GetMintable() bool
	GetOwner() sdk.AccAddress
}

// NewToken constructs a new Token instance
func NewFanToken(
	denom string,
	name string,
	maxSupply sdk.Int,
	mintable bool,
	owner sdk.AccAddress,
) FanToken {
	return FanToken{
		Denom:     denom,
		Name:      name,
		MaxSupply: maxSupply,
		Mintable:  mintable,
		Owner:     owner.String(),
	}
}

// GetDenom implements exported.TokenI
func (t FanToken) GetDenom() string {
	return t.Denom
}

// GetName implements exported.TokenI
func (t FanToken) GetName() string {
	return t.Name
}

// GetMaxSupply implements exported.TokenI
func (t FanToken) GetMaxSupply() sdk.Int {
	return t.MaxSupply
}

// GetMintable implements exported.TokenI
func (t FanToken) GetMintable() bool {
	return t.Mintable
}

// GetOwner implements exported.TokenI
func (t FanToken) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(t.Owner)
	return owner
}

func (t FanToken) String() string {
	bz, _ := yaml.Marshal(t)
	return string(bz)
}
