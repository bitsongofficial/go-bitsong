package app

import (
	"fmt"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"time"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// DefaultConfig returns a default configuration suitable for nearly all
// testing requirements.
func DefaultConfig() network.Config {
	encCfg := MakeEncodingConfig()
	chainId := "bitsong-test-1"

	return network.Config{
		Codec:             encCfg.Marshaler,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor:    NewAppConstructor(chainId),
		GenesisState:      AppModuleBasics.DefaultGenesis(encCfg.Marshaler),
		TimeoutCommit:     1 * time.Second / 2,
		ChainID:           chainId,
		NumValidators:     1,
		BondDenom:         sdk.DefaultBondDenom,
		MinGasPrices:      fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:     sdk.TokensFromConsensusPower(100000, sdk.DefaultPowerReduction),
		StakingTokens:     sdk.TokensFromConsensusPower(50000, sdk.DefaultPowerReduction),
		BondedTokens:      sdk.TokensFromConsensusPower(10000, sdk.DefaultPowerReduction),
		PruningStrategy:   "nothing",
		CleanupDir:        true,
		SigningAlgo:       string(hd.Secp256k1Type),
		KeyringOptions:    []keyring.Option{},
	}
}

func NewAppConstructor(chainId string) network.AppConstructor {
	return func(val network.ValidatorI) servertypes.Application {
		valCtx := val.GetCtx()
		appConfig := val.GetAppConfig()

		return NewBitsongApp(
			valCtx.Logger, dbm.NewMemDB(), nil, true, valCtx.Config.RootDir, EmptyAppOptions{}, EmptyWasmOpts,
			baseapp.SetMinGasPrices(appConfig.MinGasPrices), baseapp.SetChainID(chainId),
		)
	}
}
