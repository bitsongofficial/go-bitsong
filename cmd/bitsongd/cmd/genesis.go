package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

type GenesisParams struct {
	GenesisTime     time.Time
	ConsensusParams *tmtypes.ConsensusParams
	NativeCoin      []banktypes.Metadata

	StakingParams      stakingtypes.Params
	MintParams         minttypes.Params
	DistributionParams distributiontypes.Params
	GovParams          govtypes.Config
	SlashingParams     slashingtypes.Params

	CrisisConstantFee sdk.Coin

	FantokenParams fantokentypes.Params
}

func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.GenesisTime = time.Date(2021, 10, 21, 11, 0, 0, 0, time.UTC) // Oct 10, 2021 - 11:00 UTC

	genParams.NativeCoin = []banktypes.Metadata{
		{
			Description: "The native token of BitSong Network",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    appparams.MicroCoinUnit,
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    appparams.CoinUnit,
					Exponent: appparams.CoinExponent,
					Aliases:  nil,
				},
			},
			Base:    appparams.MicroCoinUnit,
			Display: appparams.CoinUnit,
		},
	}

	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 7 * 3 // 3 weeks
	genParams.StakingParams.MaxValidators = 100
	genParams.StakingParams.BondDenom = appparams.MicroCoinUnit

	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 20 * 1024 * 1024 // 20MB
	genParams.ConsensusParams.Block.MaxGas = 200_000_000        // 200.000.000 units
	genParams.ConsensusParams.Evidence.MaxAgeDuration = genParams.StakingParams.UnbondingTime
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	// genParams.ConsensusParams.Evidence.MaxAgeDuration = 172_800_000_000_000
	// genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = 100_000

	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.BlocksPerYear = 5733820
	genParams.MintParams.MintDenom = appparams.MicroCoinUnit

	genParams.DistributionParams = distributiontypes.DefaultParams()
	genParams.DistributionParams.BaseProposerReward = sdk.MustNewDecFromStr("0.01")
	genParams.DistributionParams.BonusProposerReward = sdk.MustNewDecFromStr("0.04")
	genParams.DistributionParams.CommunityTax = sdk.MustNewDecFromStr("0.02")
	genParams.DistributionParams.WithdrawAddrEnabled = true

	genParams.GovParams = govtypes.DefaultConfig()
	// genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 24 * 15 // 15 days
	// genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
	// 	appparams.MicroCoinUnit,
	// 	sdk.NewInt(512_000_000),
	// ))
	// genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.4")          // 40%
	// genParams.GovParams.TallyParams.Threshold = sdk.MustNewDecFromStr("0.5")       // 50%
	// genParams.GovParams.TallyParams.VetoThreshold = sdk.MustNewDecFromStr("0.334") // 33.40%
	// genParams.GovParams.VotingParams.VotingPeriod = time.Hour * 24 * 7             // 7 days

	genParams.SlashingParams = slashingtypes.DefaultParams()
	genParams.SlashingParams.SignedBlocksWindow = int64(10000)                       // 10000 blocks (~13.8 hr at 5 second blocks)
	genParams.SlashingParams.MinSignedPerWindow = sdk.MustNewDecFromStr("0.05")      // 5% minimum liveness
	genParams.SlashingParams.DowntimeJailDuration = time.Hour                        // 1 hour jail period
	genParams.SlashingParams.SlashFractionDoubleSign = sdk.MustNewDecFromStr("0.05") // 5% double sign slashing
	genParams.SlashingParams.SlashFractionDowntime = sdk.MustNewDecFromStr("0.01")   // 1% liveness slashing

	genParams.CrisisConstantFee = sdk.NewCoin(appparams.MicroCoinUnit, sdk.NewInt(133_333_000_000))

	genParams.FantokenParams = fantokentypes.DefaultParams()
	genParams.FantokenParams.IssueFee = sdk.NewCoin(appparams.MicroCoinUnit, sdk.NewInt(1_000_000_000))
	genParams.FantokenParams.MintFee = sdk.NewCoin(appparams.MicroCoinUnit, sdk.ZeroInt())
	genParams.FantokenParams.BurnFee = sdk.NewCoin(appparams.MicroCoinUnit, sdk.ZeroInt())

	return genParams
}

func TestnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.GenesisTime = time.Now()

	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 7 * 2 // 2 weeks

	// genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
	// 	appparams.MicroCoinUnit,
	// 	sdk.NewInt(1000000), // 1 BTSG
	// ))
	// genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.0000000001") // 0.00000001%
	// genParams.GovParams.VotingParams.VotingPeriod = time.Second * 300              // 300 seconds

	return genParams
}

func PrepareGenesis(clientCtx client.Context, appState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, genesisParams GenesisParams, chainID string) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	depCdc := clientCtx.Codec

	genDoc.ChainID = chainID
	genDoc.GenesisTime = genesisParams.GenesisTime
	genDoc.ConsensusParams = genesisParams.ConsensusParams

	stakingGenState := stakingtypes.GetGenesisStateFromAppState(depCdc, appState)
	stakingGenState.Params = genesisParams.StakingParams
	stakingGenStateBz, err := depCdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = genesisParams.MintParams
	mintGenStateBz, err := depCdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = genesisParams.DistributionParams
	distributionGenStateBz, err := depCdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// govGenState := govtypes.DefaultGenesisState()
	// govGenState.DepositParams = genesisParams.GovParams.DepositParams
	// govGenState.TallyParams = genesisParams.GovParams.TallyParams
	// govGenState.VotingParams = genesisParams.GovParams.VotingParams
	// govGenStateBz, err := cdc.MarshalJSON(govGenState)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	// }
	// appState[govtypes.ModuleName] = govGenStateBz

	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = genesisParams.CrisisConstantFee
	crisisGenStateBz, err := depCdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = genesisParams.SlashingParams
	slashingGenStateBz, err := depCdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	fantokenGenState := fantokentypes.DefaultGenesisState()
	fantokenGenState.Params = genesisParams.FantokenParams
	fantokenGenStateBz, err := depCdc.MarshalJSON(fantokenGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal fantoken genesis state: %w", err)
	}
	appState[fantokentypes.ModuleName] = fantokenGenStateBz

	return appState, genDoc, nil
}

func PrepareGenesisCmd(defaultNodeHome string, mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-genesis",
		Short: "Prepare a genesis file with initial setup",
		Long: `Prepare a genesis file with initial setup.
Examples include:
	- Setting module initial params
	- Setting denom metadata
Example:
	bitsongd prepare-genesis mainnet bitsong-2b
	- Check input genesis:
		file is at ~/.bitsongd/config/genesis.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// read genesis file
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// get genesis params
			var genesisParams GenesisParams
			network := args[0]
			if network == "testnet" {
				genesisParams = TestnetGenesisParams()
			} else if network == "mainnet" {
				genesisParams = MainnetGenesisParams()
			} else {
				return fmt.Errorf("please choose 'mainnet' or 'testnet'")
			}

			// get genesis params
			chainID := args[1]

			// run Prepare Genesis
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, genesisParams, chainID)
			if err != nil {
				return fmt.Errorf("err, %s", err.Error())
			}

			// validate genesis state
			if err = mbm.ValidateGenesis(depCdc, clientCtx.TxConfig, appState); err != nil {
				return fmt.Errorf("error validating genesis file: %s", err.Error())
			}

			// save genesis
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
