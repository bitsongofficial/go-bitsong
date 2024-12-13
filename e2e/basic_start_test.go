package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	bitsongconformance "github.com/bitsongofficial/go-bitsong/tests/e2e/conformance"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// TestBasicBitsongStart is a basic test to assert that spinning up a Bitsong network with one validator works properly.
func TestBasicBitsongStart(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Base setup
	chains := CreateThisBranchWithValsAndFullNodes(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	bitsong := chains[0].(*cosmos.CosmosChain)

	userFunds := sdkmath.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, bitsong)
	println("users", users)
	chainUser := users[0]

	bitsongconformance.ConformanceCosmWasm(t, ctx, bitsong, chainUser)

	// grpc query
	// bitsongGrpc := bitsong.GetGRPCAddress()
	// dialOpts := grpc.WithTransportCredentials(insecure.NewCredentials())
	// conn, err := grpc.NewClient(bitsongGrpc, dialOpts)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(conn)
	// defer conn.Close()
	// client := wasmtypes.NewQueryClient(conn)
	// req6 := &wasmtypes.QueryParamsRequest{}
	// resp6, err6 := client.Params(ctx, req6)
	// if err6 != nil {
	// 	log.Fatal(err)
	// }
	// require.NotNil(t, resp6)

	require.NotNil(t, ic)
	require.NotNil(t, ctx)

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
