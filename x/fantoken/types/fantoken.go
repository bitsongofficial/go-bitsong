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

// FanTokenI defines an interface for FanToken
type FanTokenI interface {
	GetSymbol() string
	GetDenom() string
	GetName() string
	GetMaxSupply() sdk.Int
	GetMintable() bool
	GetOwner() sdk.AccAddress
	GetUri() string
	GetMetaData() banktypes.Metadata
}

// NewFanToken constructs a new Token instance
func NewFanToken(
	name string,
	maxSupply sdk.Int,
	owner sdk.AccAddress,
	uri string,
	metaData banktypes.Metadata,
) FanToken {
	return FanToken{
		Name:      name,
		MaxSupply: maxSupply,
		Mintable:  true,
		Owner:     owner.String(),
		URI:       uri,
		MetaData:  metaData,
	}
}

// GetSymbol implements exported.FanTokenI
func (ft FanToken) GetSymbol() string {
	return ft.MetaData.Display
}

// GetDenom implements exported.FanTokenI
func (ft FanToken) GetDenom() string {
	return ft.MetaData.Base
}

// GetName implements exported.FanTokenI
func (ft FanToken) GetName() string {
	return ft.Name
}

// GetMaxSupply implements exported.FanTokenI
func (ft FanToken) GetMaxSupply() sdk.Int {
	return ft.MaxSupply
}

// GetMintable implements exported.FanTokenI
func (ft FanToken) GetMintable() bool {
	return ft.Mintable
}

// GetOwner implements exported.FanTokenI
func (ft FanToken) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(ft.Owner)
	return owner
}

// GetUri implements exported.FanTokenI
func (ft FanToken) GetUri() string {
	return ft.URI
}

// GetMetaData returns metadata of the fantoken
func (ft FanToken) GetMetaData() banktypes.Metadata {
	return ft.MetaData
}

func (ft FanToken) String() string {
	bz, _ := yaml.Marshal(ft)
	return string(bz)
}
