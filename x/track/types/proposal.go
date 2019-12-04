package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeTrackVerify defines the type for a TrackVerifyProposal
	ProposalTypeTrackVerify = "TrackVerify"
)

// Assert TrackVerifyProposal implements govtypes.Content at compile-time
var _ govtypes.Content = TrackVerifyProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeTrackVerify)
	govtypes.RegisterProposalTypeCodec(TrackVerifyProposal{}, "go-bitsong/TrackVerifyProposal")
}

// TrackVerifyProposal verify a track
type TrackVerifyProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	TrackID     uint64 `json:"id" yaml:"id"`
}

// NewArtistVerifyProposal creates a new artist verify proposal.
func NewTrackVerifyProposal(title, description string, id uint64) TrackVerifyProposal {
	return TrackVerifyProposal{title, description, id}
}

// GetTitle returns the title of a community pool spend proposal.
func (vp TrackVerifyProposal) GetTitle() string { return vp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (vp TrackVerifyProposal) GetDescription() string { return vp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (vp TrackVerifyProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (vp TrackVerifyProposal) ProposalType() string { return ProposalTypeTrackVerify }

// ValidateBasic runs basic stateless validity checks
func (vp TrackVerifyProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, vp)
	if err != nil {
		return err
	}

	// TODO:
	// Only owner can open proposal?
	if vp.TrackID == 0 {
		return ErrUnknownTrack(DefaultCodespace, "unknown track")
	}

	return nil
}

// String implements the Stringer interface.
func (vp TrackVerifyProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Track Verify Proposal:
  Title:       %s
  Description: %s
  Track ID: %d
`, vp.Title, vp.Description, vp.TrackID))
	return b.String()
}
