package e2e

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
)

const (
	v019UpgradeName = "v019"
)

func TestBasicBitsongSlashing(t *testing.T) {
	repo, version := GetDockerImageInfo()
	V019PatchTest(t, chainName, version, repo, upgradeName)
}

func V019PatchTest(t *testing.T, chainName, upgradeBranchVersion, upgradeRepo, upgradeName string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	t.Log(chainName, upgradeBranchVersion, upgradeRepo, upgradeName)

	// expect validator to participate in each block
	previousVersionGenesis := []cosmos.GenesisKV{
		{Key: "app_state.gov.params.voting_period", Value: VotingPeriod},
		{Key: "app_state.gov.params.max_deposit_period", Value: MaxDepositPeriod},
		{Key: "app_state.gov.params.min_deposit.0.denom", Value: Denom},
		{Key: "app_state.slashing.params.signed_blocks_window", Value: "1"},
		{Key: "app_state.slashing.params.min_signed_per_window", Value: "1.000000000000000000"},
	}

	cfg := baseCfg
	cfg.ModifyGenesis = cosmos.ModifyGenesis(previousVersionGenesis)
	cfg.Images = []ibc.DockerImage{baseChain}

	numVals, numNodes := 4, 1
	chains := CreateICTestBitsongChainCustomConfig(t, numVals, numNodes, cfg)
	chain := chains[0].(*cosmos.CosmosChain)

	ic, ctx, client, _ := BuildInitialChain(t, chains)
	t.Cleanup(func() {
		_ = ic.Close()
	})
	delegatorFunds := sdkmath.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), delegatorFunds, chain)
	users2, err := interchaintest.GetAndFundTestUserWithMnemonic(
		ctx,
		t.Name(),
		"clap lounge repair barely exile forward kangaroo festival staff stove mimic reveal sudden mosquito coral antique inch kite impact track maple mean coffee slab",
		delegatorFunds,
		chain,
	)
	require.NoError(t, err)
	btsgDelegator := users[0]

	// get validator addresses
	vals, err := QueryAllValidators(t, chain, ctx)
	require.NoError(t, err)
	valAddr := vals.Validators[0].OperatorAddress

	for _, val := range vals.Validators[1:3] {
		// for the first validator delegate only 1000
		amnt := "10000ubtsg"

		delegate := []string{
			chain.Config().Bin, "tx", "staking", "delegate",
			val.OperatorAddress,
			amnt,
			"--from", btsgDelegator.KeyName(),
			"--chain-id", chain.Config().ChainID,
			"--home", chain.HomeDir(),
			"--node", chain.GetRPCAddress(),
			"--keyring-backend", keyring.BackendTest,
			"-y",
		}
		_, _, err = chain.Exec(ctx, delegate, nil)
		require.NoError(t, err)
	}

	// delegate to first validator with fresh account
	delegate := []string{
		chain.Config().Bin, "tx", "staking", "delegate",
		valAddr,
		"1000ubtsg",
		"--from", users2.KeyName(),
		"--chain-id", chain.Config().ChainID,
		"--home", chain.HomeDir(),
		"--node", chain.GetRPCAddress(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}
	_, _, err = chain.Exec(ctx, delegate, nil)
	require.NoError(t, err)

	// TODO: get to work so validator from this node gets slashed offline
	// stop validator from participating in voting
	chain.Validators[0].PauseContainer(ctx)
	_ = testutil.WaitForBlocks(ctx, int(2), chain)
	chain.Validators[0].StartContainer(ctx)

	// query slashing rewards
	slashEvents, err := QuerySlashedEvents(t, ctx, chain)
	require.NoError(t, err)
	fmt.Println("SLASHING EVENTS:", slashEvents.Slashes)

	// query rewards for both delegators
	res, err := QueryStakingDistributionRewards(t, chain, ctx, btsgDelegator.FormattedAddress(), valAddr)
	require.NoError(t, err)
	fmt.Println("GOOD RESPONSE:", res.Rewards)

	res, err = QueryStakingDistributionRewards(t, chain, ctx, users2.FormattedAddress(), valAddr)
	fmt.Println("BAD RESPONSE:", res)
	if err != nil {
		fmt.Println("BAD RESPONSE ERROR: ", err)
	}

	// upgrade to v019 to apply patch
	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	haltHeight := height + haltHeightDelta
	proposalID := SubmitUpgradeProposal(t, ctx, chain, btsgDelegator, upgradeName, haltHeight, users2)

	proposalIDInt, err := strconv.ParseInt(proposalID, 10, 64)
	require.NoError(t, err, "failed to parse proposal ID")

	ValidatorVoting(t, ctx, chain, proposalIDInt, height, haltHeight)
	UpgradeNodes(t, ctx, chain, client, haltHeight, upgradeRepo, upgradeBranchVersion)

	// query rewards after upgrade
	res, err = QueryStakingDistributionRewards(t, chain, ctx, btsgDelegator.FormattedAddress(), valAddr)
	fmt.Println("GOOD RESPONSE:", res)
	require.NoError(t, err, "failed to patch")

	// ensure reward query is resolved
	res, err = QueryStakingDistributionRewards(t, chain, ctx, users2.FormattedAddress(), valAddr)
	fmt.Println("GOOD RESPONSE:", res)
	require.NoError(t, err, "failed to patch ")
}

// QueryStakingDistributionRewards queries the rewards for a delegator
func QueryStakingDistributionRewards(t *testing.T, c *cosmos.CosmosChain, ctx context.Context, delegator string, validator string) (*distributiontypes.QueryDelegationRewardsResponse, error) {
	res, err := distributiontypes.NewQueryClient(c.GetNode().GrpcConn).
		DelegationRewards(ctx, &distributiontypes.QueryDelegationRewardsRequest{
			DelegatorAddress: delegator,
			ValidatorAddress: validator,
		})
	return res, err
}

// QueryAllValidators queries the rewards for a delegator
func QueryAllValidators(t *testing.T, c *cosmos.CosmosChain, ctx context.Context) (*stakingtypes.QueryValidatorsResponse, error) {
	res, err := stakingtypes.NewQueryClient(c.GetNode().GrpcConn).
		Validators(ctx, &stakingtypes.QueryValidatorsRequest{})
	return res, err
}

// QueryAllValidators queries the rewards for a delegator
func QuerySlashedEvents(t *testing.T, ctx context.Context, c *cosmos.CosmosChain) (*distributiontypes.QueryValidatorSlashesResponse, error) {
	res, err := distributiontypes.NewQueryClient(c.GetNode().GrpcConn).
		ValidatorSlashes(ctx, &distributiontypes.QueryValidatorSlashesRequest{})
	return res, err
}
