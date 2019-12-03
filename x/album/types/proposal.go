package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeAlbumVerify defines the type for a AlbumVerifyProposal
	ProposalTypeAlbumVerify = "AlbumVerify"
)

// Assert AlbumVerifyProposal implements govtypes.Content at compile-time
var _ govtypes.Content = AlbumVerifyProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAlbumVerify)
	govtypes.RegisterProposalTypeCodec(AlbumVerifyProposal{}, "go-bitsong/AlbumVerifyProposal")
}

// AlbumVerifyProposal verify an album
type AlbumVerifyProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	AlbumID     uint64 `json:"id" yaml:"id"`
}

// NewAlbumVerifyProposal creates a new album verify proposal.
func NewAlbumVerifyProposal(title, description string, id uint64) AlbumVerifyProposal {
	return AlbumVerifyProposal{title, description, id}
}

// GetTitle returns the title of a community pool spend proposal.
func (vp AlbumVerifyProposal) GetTitle() string { return vp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (vp AlbumVerifyProposal) GetDescription() string { return vp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (vp AlbumVerifyProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (vp AlbumVerifyProposal) ProposalType() string { return ProposalTypeAlbumVerify }

// ValidateBasic runs basic stateless validity checks
func (vp AlbumVerifyProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, vp)
	if err != nil {
		return err
	}

	// TODO:
	// Only owner can open proposal?
	if vp.AlbumID == 0 {
		return ErrUnknownAlbum(DefaultCodespace, "unknown album")
	}

	return nil
}

// String implements the Stringer interface.
func (vp AlbumVerifyProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Album Verify Proposal:
  Title:       %s
  Description: %s
  AlbumID: %d
`, vp.Title, vp.Description, vp.AlbumID))
	return b.String()
}
