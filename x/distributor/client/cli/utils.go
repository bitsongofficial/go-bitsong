package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io/ioutil"
)

type (
	DistributorVerifyProposalJSON struct {
		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Address     sdk.AccAddress `json:"address" yaml:"address"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)

func ParseDistributorVerifyProposalJSON(cdc *codec.Codec, proposalFile string) (DistributorVerifyProposalJSON, error) {
	proposal := DistributorVerifyProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
