package types

import (
	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const FanTokenDecimal = 6

var (
	_ proto.Message = &FanToken{}
)

// TokenI defines an interface for Token
type FanTokenI interface {
	GetSymbol() string
	GetDenom() string
	GetName() string
	GetMaxSupply() sdk.Int
	GetMintable() bool
	GetOwner() sdk.AccAddress
	GetMetaData() banktypes.Metadata
}

// NewToken constructs a new Token instance
func NewFanToken(
	name string,
	maxSupply sdk.Int,
	owner sdk.AccAddress,
	metaData banktypes.Metadata,
) FanToken {
	return FanToken{
		Name:      name,
		MaxSupply: maxSupply,
		Mintable:  true,
		Owner:     owner.String(),
		MetaData:  metaData,
	}
}

// GetSymbol implements exported.TokenI
func (t FanToken) GetSymbol() string {
	return t.MetaData.Display
}

// GetDenom implements exported.TokenI
func (t FanToken) GetDenom() string {
	return t.MetaData.Base
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

// GetMetaData returns metadata of the fantoken
func (t FanToken) GetMetaData() banktypes.Metadata {
	return t.MetaData
}

func (t FanToken) String() string {
	bz, _ := yaml.Marshal(t)
	return string(bz)
}
