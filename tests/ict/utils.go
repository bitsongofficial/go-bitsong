package e2e

import (
	"context"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

func sendTokens(ctx context.Context, chain *cosmos.CosmosChain, from, to ibc.Wallet, token string, amount int64) (ibc.WalletAmount, error) {
	if token == "" {
		token = chain.Config().Denom
	}

	sendAmt := ibc.WalletAmount{
		Address: to.FormattedAddress(),
		Denom:   token,
		Amount:  math.NewInt(amount),
	}
	err := chain.SendFunds(ctx, from.KeyName(), sendAmt)
	return sendAmt, err
}
