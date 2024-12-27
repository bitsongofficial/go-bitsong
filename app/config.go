package app

import (
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/bitsongofficial/go-bitsong/app/params"
	dbm "github.com/cometbft/cometbft-db"
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

	return network.Config{
		Codec:             encCfg.Marshaler,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor:    NewAppConstructor(encCfg),
		GenesisState:      AppModuleBasics.DefaultGenesis(encCfg.Marshaler),
		TimeoutCommit:     1 * time.Second / 2,
		ChainID:           "bitsong-test-1",
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

func NewAppConstructor(encodingCfg params.EncodingConfig, opts ...wasmkeeper.Option) network.AppConstructor {
	return func(val network.ValidatorI) servertypes.Application {
		return NewBitsongApp(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true, EmptyAppOptions{}, opts,
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
		)
	}
}
