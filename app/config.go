package app

// import (
// 	"fmt"

// 	"github.com/bitsongofficial/go-bitsong/app/params"
// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	"github.com/cosmos/cosmos-sdk/crypto/hd"
// 	"github.com/cosmos/cosmos-sdk/crypto/keyring"
// 	servertypes "github.com/cosmos/cosmos-sdk/server/types"

// 	"time"

// 	storetypes "github.com/cosmos/cosmos-sdk/store/types"
// 	"github.com/cosmos/cosmos-sdk/testutil/network"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// 	dbm "github.com/cometbft/cometbft-db"
// )

// DefaultConfig returns a default configuration suitable for nearly all
// testing requirements.
// func DefaultConfig() network.Config {
// 	encCfg := MakeEncodingConfig()

// 	return network.Config{
// 		Codec:             encCfg.Marshaler,
// 		TxConfig:          encCfg.TxConfig,
// 		LegacyAmino:       encCfg.Amino,
// 		InterfaceRegistry: encCfg.InterfaceRegistry,
// 		AccountRetriever:  authtypes.AccountRetriever{},
// 		AppConstructor:    NewAppConstructor(encCfg),
// 		GenesisState:      ModuleBasics.DefaultGenesis(encCfg.Marshaler),
// 		TimeoutCommit:     1 * time.Second / 2,
// 		ChainID:           "bitsong-test-1",
// 		NumValidators:     1,
// 		BondDenom:         sdk.DefaultBondDenom,
// 		MinGasPrices:      fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
// 		AccountTokens:     sdk.TokensFromConsensusPower(100000, sdk.DefaultPowerReduction),
// 		StakingTokens:     sdk.TokensFromConsensusPower(50000, sdk.DefaultPowerReduction),
// 		BondedTokens:      sdk.TokensFromConsensusPower(10000, sdk.DefaultPowerReduction),
// 		PruningStrategy:   storetypes.PruningOptionNothing,
// 		CleanupDir:        true,
// 		SigningAlgo:       string(hd.Secp256k1Type),
// 		KeyringOptions:    []keyring.Option{},
// 	}
// }

// func NewAppConstructor(encodingCfg params.EncodingConfig) network.AppConstructor {
// 	return func(val network.Validator) servertypes.Application {
// 		return NewBitsongApp(
// 			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
// 			encodingCfg,
// 			simapp.EmptyAppOptions{},
// 			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
// 		)
// 	}
// }
