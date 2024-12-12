package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cosmossdk.io/math"
	tmcfg "github.com/cometbft/cometbft/config"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

// calculate rewards for each delegator
type DifferingDelegation struct {
	DelegatorAddress string
	ValidatorAddress string
	ExpectedStake    math.LegacyDec
	ActualStake      math.LegacyDec
}

type v019LogicParams struct {
	StateFile string
}

// InitFromStateCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func V019(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v019 state_exported.json",
		Short: "calculates all validators delegators distribution errors.",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			stateFile := args[0]

			// attempt to lookup address from Keybase if no address was provided
			kb, err := keyring.New(sdk.KeyringServiceName(), "test", clientCtx.HomeDir, bufio.NewReader(cmd.InOrStdin()), clientCtx.Codec)
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}
			clientCtx.Keyring = kb

			genParams := v019LogicParams{
				StateFile: stateFile,
			}

			return v019Logic(config, clientCtx, genParams)
		},
	}

	return cmd
}

func v019Logic(_ *tmcfg.Config, cliCtx client.Context, genParams v019LogicParams) error {
	differingDelegations, err := V019ImplLogic(cliCtx, genParams)
	if err != nil {
		return fmt.Errorf("failed, dang: %w", err)
	}

	if err := printDifferingDelegations(differingDelegations); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Verification Complete")
	return nil
}

func V019ImplLogic(clientCtx client.Context, params v019LogicParams) (*[]DifferingDelegation, error) {
	if !tmos.FileExists(params.StateFile) {
		return nil, fmt.Errorf("%s does not exist", params.StateFile)
	}

	// print state export params
	fmt.Println(params.String())

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	stateBz, err := os.ReadFile(params.StateFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read state export file: %w", err)
	}
	genDoc := tmtypes.GenesisDoc{}
	err = tmjson.Unmarshal(stateBz, &genDoc)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal state export file: %w", err)
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	appState, err := genutiltypes.GenesisStateFromGenDoc(genDoc)
	if err != nil {
		return nil, err
	}

	var staking stakingtypes.GenesisState
	var distribution distrtypes.GenesisState
	var slashing slashingtypes.GenesisState
	var bank banktypes.GenesisState
	var mint minttypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appState[stakingtypes.ModuleName], &staking)
	clientCtx.Codec.MustUnmarshalJSON(appState[distrtypes.ModuleName], &distribution)
	clientCtx.Codec.MustUnmarshalJSON(appState[banktypes.ModuleName], &bank)
	clientCtx.Codec.MustUnmarshalJSON(appState[minttypes.ModuleName], &mint)
	clientCtx.Codec.MustUnmarshalJSON(appState[slashingtypes.ModuleName], &slashing)

	// get delegations from x/staking
	delByValAddr := make(map[string][]stakingtypes.Delegation)

	for _, del := range staking.Delegations {
		valAddr := del.ValidatorAddress
		if _, ok := delByValAddr[valAddr]; !ok {
			delByValAddr[valAddr] = make([]stakingtypes.Delegation, 0)
		}
		delByValAddr[valAddr] = append(delByValAddr[valAddr], del)
	}

	var differingDelegations []DifferingDelegation

	for _, delStartInfo := range distribution.DelegatorStartingInfos {
		for _, del := range delByValAddr[delStartInfo.ValidatorAddress] {
			if delStartInfo.DelegatorAddress == del.DelegatorAddress {
				for _, validator := range staking.GetValidators() {
					if validator.OperatorAddress == delStartInfo.ValidatorAddress {
						delShares := del.GetShares()
						startingInfoStake := delStartInfo.StartingInfo.Stake
						// startingPeriod := delStartInfo.StartingInfo.PreviousPeriod
						delegationShares := validator.DelegatorShares
						valTotalStakedTokens := validator.Tokens
						startingHeight := delStartInfo.StartingInfo.Height
						endingHeight := uint64(genDoc.InitialHeight)
						if endingHeight > startingHeight {
							// 1. get slash events for validator between start & end+1
							// for i, event := range distribution.ValidatorSlashEvents {

							// }
							// 2. reimplement applying slash events
							// k.IterateValidatorSlashEventsBetween(ctx, del.GetValidatorAddr(), startingHeight, endingHeight,
							// 	func(height uint64, event types.ValidatorSlashEvent) (stop bool) {
							// 		endingPeriod := event.ValidatorPeriod
							// 		if endingPeriod > startingPeriod {
							// 			rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)
							// 			startingInfoStake = stake.MulTruncate(math.LegacyOneDec().Sub(event.Fraction))
							// 			startingPeriod = endingPeriod
							// 		}
							// 		return false
							// 	},
							// )
						}

						currentStake := CustomTokensFromShares(valTotalStakedTokens, delegationShares, delShares)
						if startingInfoStake.GT(currentStake) {
							// create list of delegators and validators that calculations differ than expected
							differingDelegations = append(differingDelegations, DifferingDelegation{
								DelegatorAddress: del.DelegatorAddress,
								ValidatorAddress: validator.OperatorAddress,
								ExpectedStake:    startingInfoStake,
								ActualStake:      currentStake,
							})
						}
					}
				}
			}
		}
	}

	return &differingDelegations, nil
}

// t = (shares * tokens) / total_delegators
func CustomTokensFromShares(v math.Int, ds math.LegacyDec, shares sdk.Dec) math.LegacyDec {
	return (shares.MulInt(v)).Quo(ds)
}

func (s *v019LogicParams) String() string {
	return fmt.Sprintf(` 
  State File: %s
`, s.StateFile)
}

func (s *v019LogicParams) Validate() error {
	if s.StateFile == "" {
		return fmt.Errorf("state file cannot be empty")
	}

	return nil
}

func printDifferingDelegations(differingDelegations *[]DifferingDelegation) error {
	filename := "differing_delegations.json"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	json.NewEncoder(f).Encode(differingDelegations)
	return nil
}
