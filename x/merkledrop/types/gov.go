package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const ProposalTypeUpdateFees = "UpdateMerkledropFeesProposal"

func init() {
	v1beta1.RegisterProposalType(ProposalTypeUpdateFees)
	// v1beta1.RegisterProposalTypeCodec(&UpdateFeesProposal{}, "go-bitsong/merkledrop/UpdateFeesProposal")
}

var _ v1beta1.Content = &UpdateFeesProposal{}

func NewUpdateFeesProposal(title, description string, creationFee sdk.Coin) v1beta1.Content {
	return &UpdateFeesProposal{
		Title:       title,
		Description: description,
		CreationFee: creationFee,
	}
}

func (p *UpdateFeesProposal) GetTitle() string { return p.Title }

func (p *UpdateFeesProposal) GetDescription() string { return p.Description }

func (p *UpdateFeesProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateFeesProposal) ProposalType() string { return ProposalTypeUpdateFees }

func (p *UpdateFeesProposal) ValidateBasic() error {
	err := v1beta1.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if err := p.CreationFee.Validate(); err != nil {
		return err
	}

	return nil
}

func (p UpdateFeesProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Merkledrop Fees Proposal:
  Title:       %s
  Description: %s
  Creation Fee:   %s
`, p.Title, p.Description, p.CreationFee))
	return b.String()
}
