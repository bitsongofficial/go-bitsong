package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Constants pertaining to a Meta object
const (
	MaxNameLength int = 140
)

// Meta defines an interface that an artist must implement. It contains
// information such as the title along with the routing
// information for the appropriate handler to process the artist. Content can
// have additional fields, which will handled by a artist's Handler.
type Meta interface {
	GetName() string
	ArtistRoute() string
	ValidateBasic() sdk.Error
	String() string
}

// ValidateAbstract validates an artist's abstract meta returning an error
// if invalid.
func ValidateAbstract(codespace sdk.CodespaceType, m Meta) sdk.Error {
	name := m.GetName()
	if len(strings.TrimSpace(name)) == 0 {
		return ErrInvalidArtistMeta(codespace, "artist name cannot be blank")
	}
	if len(name) > MaxNameLength {
		return ErrInvalidArtistMeta(codespace, fmt.Sprintf("artist name is longer than max length of %d", MaxNameLength))
	}

	return nil
}

// General Meta
type GeneralMeta struct {
	Name string `json:"name" yaml:"name"`
}

func NewGeneralMeta(name string) Meta {
	return GeneralMeta{name}
}

// Implements Artist Interface
var _ Meta = GeneralMeta{}

// nolint
func (gm GeneralMeta) GetName() string          { return gm.Name }
func (gm GeneralMeta) ArtistRoute() string      { return RouterKey }
func (gm GeneralMeta) ValidateBasic() sdk.Error { return ValidateAbstract(DefaultCodespace, gm) }

func (gm GeneralMeta) String() string {
	return fmt.Sprintf(`General Meta:
  Name:       %s
`, gm.Name)
}
