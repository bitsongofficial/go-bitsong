package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bitsongofficial/go-bitsong/app/params"
	cosmosdb "github.com/cosmos/cosmos-db"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cosmossdk.io/log"
	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/libs/bytes"
	tmcli "github.com/cometbft/cometbft/libs/cli"

	// rosettaCmd "cosmossdk.io/tools/rosetta/cmd"

	bitsong "github.com/bitsongofficial/go-bitsong/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
)

// NewRootCmd creates a new root command for bitsongd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := bitsong.MakeEncodingConfig()

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(bitsong.Bech32PrefixAccAddr, bitsong.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(bitsong.Bech32PrefixValAddr, bitsong.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(bitsong.Bech32PrefixConsAddr, bitsong.Bech32PrefixConsPub)
	cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())
	cfg.Seal()

	appOptions := make(simtestutil.AppOptionsMap, 0)

	tempDir := tempDir()
	// cleanup temp dir after we are done with the tempApp, so we don't leave behind a
	// new temporary directory for every invocation. See https://github.com/CosmWasm/wasmd/issues/2017
	defer os.RemoveAll(tempDir)

	tempApp := bitsong.NewBitsongApp(
		log.NewNopLogger(),
		cosmosdb.NewMemDB(),
		nil,
		true,
		tempDir,
		appOptions,
		[]wasmkeeper.Option{},
	)

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.FlagBroadcastMode).
		WithHomeDir(bitsong.DefaultNodeHome).
		WithViper("") // In bitsong, we don't use any prefix for env variables.

	// Allows you to add extra params to your client.toml
	// gas, gas-price, gas-adjustment, fees, note, etc.
	SetCustomEnvVariablesFromClientToml(initClientCtx)

	rootCmd := &cobra.Command{
		Use:   version.AppName,
		Short: "Bitsong Network Application",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := tmcfg.DefaultConfig()
			// add customizations to tendermint configuration here
			// customTMConfig.P2P.MaxNumInboundPeers = 100
			// customTMConfig.P2P.MaxNumOutboundPeers = 40

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	if err := autoCliOpts(initClientCtx, tempApp).EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

// tempDir create a temporary directory to initialize the command line client
func tempDir() string {
	dir, err := os.MkdirTemp("", "bitsongd")
	if err != nil {
		panic(fmt.Sprintf("failed creating temp directory: %s", err.Error()))
	}

	return dir
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.

	type CustomAppConfig struct {
		serverconfig.Config

		Wasm wasmtypes.WasmConfig `mapstructure:"wasm"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In simapp, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0.001ubtsg"
	srvCfg.API.Enable = true
	srvCfg.API.Swagger = true
	// srvCfg.BaseConfig.IAVLDisableFastNode = true // disable fastnode by default

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Wasm:   wasmtypes.DefaultWasmConfig(),
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate +
		wasmtypes.DefaultConfigTemplate()

	return customAppTemplate, customAppConfig
}

// Reads the custom extra values in the config.toml file if set.
// If they are, then use them.
func SetCustomEnvVariablesFromClientToml(ctx client.Context) {
	configFilePath := filepath.Join(ctx.HomeDir, "config", "client.toml")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return
	}

	viper := ctx.Viper
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	setEnvFromConfig := func(key string, envVar string) {
		// if the user sets the env key manually, then we don't want to override it
		if os.Getenv(envVar) != "" {
			return
		}

		// reads from the config file
		val := viper.GetString(key)
		if val != "" {
			// Sets the env for this instance of the app only.
			os.Setenv(envVar, val)
		}
	}

	// gas
	setEnvFromConfig("gas", "BITSONGD_GAS")
	setEnvFromConfig("gas-prices", "BITSONGD_GAS_PRICES")
	setEnvFromConfig("gas-adjustment", "BITSONGD_GAS_ADJUSTMENT")
	// fees
	setEnvFromConfig("fees", "BITSONGD_FEES")
	setEnvFromConfig("fee-account", "BITSONGD_FEE_ACCOUNT")
	// memo
	setEnvFromConfig("note", "BITSONGD_NOTE")
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {

	rootCmd.AddCommand(
		genutilcli.InitCmd(bitsong.AppModuleBasics, bitsong.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		ConfigCmd(),
		pruning.Cmd(newApp, bitsong.DefaultNodeHome),
	)

	server.AddCommands(rootCmd, bitsong.DefaultNodeHome, newApp, appExport, addModuleInitFlags)
	server.AddTestnetCreatorCommand(rootCmd, newTestnetApp, addModuleInitFlags)
	wasmcli.ExtendUnsafeResetAllCmd(rootCmd)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		PrepareGenesisCmd(bitsong.DefaultNodeHome, bitsong.AppModuleBasics),
		genesisCommand(encodingConfig),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	wasm.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(encodingConfig params.EncodingConfig, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.GenesisCoreCommand(encodingConfig.TxConfig, bitsong.AppModuleBasics, bitsong.DefaultNodeHome)
	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	bitsong.AppModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// type appCreator struct {
// 	encCfg params.EncodingConfig
// }

// newApp creates the application
func newApp(
	logger log.Logger,
	db cosmosdb.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	loadLatest := true

	baseappOptions := server.DefaultBaseappOptions(appOpts)

	return bitsong.NewBitsongApp(
		logger,
		db,
		traceStore,
		loadLatest,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		appOpts,
		wasmOpts,
		baseappOptions...,
	)
}

// newTestnetApp starts by running the normal newApp method. From there, the app interface returned is modified in order
// for a testnet to be created from the provided app.
func newTestnetApp(logger log.Logger, db cosmosdb.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	// Create an app and type cast to an BitsongApp
	app := newApp(logger, db, traceStore, appOpts)
	bitsongApp, ok := app.(*bitsong.BitsongApp)
	if !ok {
		panic("app created from newApp is not of type bitsongApp")
	}

	newValAddr, ok := appOpts.Get(server.KeyNewValAddr).(bytes.HexBytes)
	if !ok {
		panic("newValAddr is not of type bytes.HexBytes")
	}
	newValPubKey, ok := appOpts.Get(server.KeyUserPubKey).(crypto.PubKey)
	if !ok {
		panic("newValPubKey is not of type crypto.PubKey")
	}
	newOperatorAddress, ok := appOpts.Get(server.KeyNewOpAddr).(string)
	if !ok {
		panic("newOperatorAddress is not of type string")
	}
	upgradeToTrigger, ok := appOpts.Get(server.KeyTriggerTestnetUpgrade).(string)
	if !ok {
		panic("upgradeToTrigger is not of type string")
	}
	// Make modifications to the normal BitsongApp required to run the network locally
	return bitsong.InitBitsongAppForTestnet(bitsongApp, newValAddr, newValPubKey, newOperatorAddress, upgradeToTrigger)

}

// appExport creates a new wasm app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db cosmosdb.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var wasmApp *bitsong.BitsongApp
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	var emptyWasmOpts []wasmkeeper.Option
	wasmApp = bitsong.NewBitsongApp(
		logger,
		db,
		traceStore,
		height == -1,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		appOpts,
		emptyWasmOpts,
	)

	if height != -1 {
		if err := wasmApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return wasmApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
func autoCliOpts(initClientCtx client.Context, tempApp *bitsong.BitsongApp) autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range tempApp.ModuleManager().Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(tempApp.ModuleManager().Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
		ClientCtx:             initClientCtx}
}
