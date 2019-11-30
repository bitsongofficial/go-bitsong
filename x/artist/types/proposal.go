package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeVerifyArtist defines the type for a ArtistVerifyProposal
	ProposalTypeVerifyArtist = "VerifyArtist"
)

// Assert ArtistVerifyProposal implements govtypes.Content at compile-time
var _ govtypes.Content = ArtistVerifyProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeVerifyArtist)
	govtypes.RegisterProposalTypeCodec(ArtistVerifyProposal{}, "go-bitsong/ArtistVerifyProposal")
}

// ArtistVerifyProposal verify an artist profile
type ArtistVerifyProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	ArtistID    uint64 `json:"id" yaml:"id"`
}

// NewArtistVerifyProposal creates a new artist verify proposal.
func NewArtistVerifyProposal(title, description string, id uint64) ArtistVerifyProposal {
	return ArtistVerifyProposal{title, description, id}
}

// GetTitle returns the title of a community pool spend proposal.
func (vp ArtistVerifyProposal) GetTitle() string { return vp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (vp ArtistVerifyProposal) GetDescription() string { return vp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (vp ArtistVerifyProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (vp ArtistVerifyProposal) ProposalType() string { return ProposalTypeVerifyArtist }

// ValidateBasic runs basic stateless validity checks
func (vp ArtistVerifyProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, vp)
	if err != nil {
		return err
	}

	// TODO:
	// Only owner can open proposal?
	if vp.ArtistID == 0 {
		return ErrUnknownArtist(DefaultCodespace, "unknown artist")
	}

	return nil
}

// String implements the Stringer interface.
func (vp ArtistVerifyProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Artist Verify Proposal:
  Title:       %s
  Description: %s
`, vp.Title, vp.Description))
	return b.String()
}
