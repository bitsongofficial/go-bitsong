package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmcli "github.com/CosmWasm/wasmd/x/wasm/client/cli"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dbm "github.com/cometbft/cometbft-db"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/libs/log"

	rosettaCmd "cosmossdk.io/tools/rosetta/cmd"

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

	return rootCmd, encodingConfig
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
	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(bitsong.ModuleBasics, bitsong.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		ConfigCmd(),
		pruning.PruningCmd(ac.newApp),
		// NewTestnetCmd(bitsong.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		// AddGenesisIcaCmd(bitsong.DefaultNodeHome),
		// DebugCmd(),
	)

	server.AddCommands(rootCmd, bitsong.DefaultNodeHome, ac.newApp, ac.appExport, addModuleInitFlags)
	wasmcli.ExtendUnsafeResetAllCmd(rootCmd)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		PrepareGenesisCmd(bitsong.DefaultNodeHome, bitsong.ModuleBasics),
		genesisCommand(encodingConfig),
		InitFromStateCmd(bitsong.DefaultNodeHome),
		VerifySlashedDelegatorsV018(bitsong.DefaultNodeHome),
		queryCommand(),
		txCommand(),
		keys.Commands(bitsong.DefaultNodeHome),
	)
	// add rosetta
	rootCmd.AddCommand(rosettaCmd.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Marshaler))
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	wasm.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(encodingConfig params.EncodingConfig, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.GenesisCoreCommand(encodingConfig.TxConfig, bitsong.ModuleBasics, bitsong.DefaultNodeHome)
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
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	bitsong.ModuleBasics.AddQueryCommands(cmd)
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
		authcmd.GetAuxToFeeCommand(),
	)

	bitsong.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

// newApp creates the application
func (ac appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
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
		appOpts,
		wasmOpts,
		baseappOptions...,
	)
}

// appExport creates a new wasm app (optionally at a given height) and exports state.
func (ac appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
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
