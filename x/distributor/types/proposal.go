package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strings"
)

const (
	ProposalTypeDistributorVerify = "DistributorVerify"
)

var _ govtypes.Content = DistributorVerifyProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeDistributorVerify)
	govtypes.RegisterProposalTypeCodec(DistributorVerifyProposal{}, "go-bitsong/DistributorVerifyProposal")
}

type DistributorVerifyProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Address     sdk.AccAddress `json:"address" yaml:"address"`
}

func NewDistributorVerifyProposal(title, description string, accAddr sdk.AccAddress) DistributorVerifyProposal {
	return DistributorVerifyProposal{title, description, accAddr}
}

// GetTitle returns the title of a community pool spend proposal.
func (vp DistributorVerifyProposal) GetTitle() string { return vp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (vp DistributorVerifyProposal) GetDescription() string { return vp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (vp DistributorVerifyProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (vp DistributorVerifyProposal) ProposalType() string { return ProposalTypeDistributorVerify }

// ValidateBasic runs basic stateless validity checks
func (vp DistributorVerifyProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, vp)
	if err != nil {
		return err
	}

	// TODO:
	// Add more check

	return nil
}

// String implements the Stringer interface.
func (vp DistributorVerifyProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Distributor Verify Proposal:
  Title:       %s
  Description: %s
`, vp.Title, vp.Description))
	return b.String()
}
