package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	bitsongconformance "github.com/bitsongofficial/go-bitsong/tests/e2e/conformance"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// TestBasicBtsgStart is a basic test to assert that spinning up a Bitsong network with one validator works properly.
func TestBasicBtsgStart(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Base setup
	chains := CreateThisBranchChain(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	chain := chains[0].(*cosmos.CosmosChain)

	userFunds := sdkmath.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	println("users", users)
	chainUser := users[0]

	bitsongconformance.ConformanceCosmWasm(t, ctx, chain, chainUser)

	require.NotNil(t, ic)
	require.NotNil(t, ctx)

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
