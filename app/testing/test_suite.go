package apptesting

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakinghelper "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stretchr/testify/suite"
)

var (
	SecondaryDenom       = "uakt"
	SecondaryAmount      = math.NewInt(100000000)
	baseTestAccts        = []sdk.AccAddress{}
	defaultTestStartTime = time.Now().UTC()
	testDescription      = stakingtypes.NewDescription("test_moniker", "test_identity", "test_website", "test_security_contact", "test_details")
)

type KeeperTestHelper struct {
	suite.Suite

	App         *app.BitsongApp                 // Mock bitsong application
	Ctx         sdk.Context                     // simulated context
	QueryHelper *baseapp.QueryServiceTestHelper // GRPC query simulator
	TestAccs    []sdk.AccAddress                // Test accounts

	StakingHelper *stakinghelper.Helper // Useful staking helpers

	// set to true if any method alters baseapp/abci is used.
	// controls whether or not to reuse app instance or set new one.
	hasUsedAbci bool
	withCaching bool
}

func init() {
	baseTestAccts = CreateRandomAccounts(3)
}

func (s *KeeperTestHelper) Reset() {
	if s.hasUsedAbci || !s.withCaching {
		s.withCaching = true
		s.Setup()
	} else {
		s.setupGeneral()
	}
}

func (s *KeeperTestHelper) Setup() {
	dir, err := os.MkdirTemp("", "bitsongd-test-home")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}

	// Create minimal config files for testing
	err = s.createTestConfigFiles(dir)
	if err != nil {
		panic(fmt.Sprintf("failed creating test config files: %v", err))
	}

	s.T().Cleanup(func() { os.RemoveAll(dir); s.withCaching = false })
	s.App = app.SetupWithCustomHome(false, dir)
	// configure ctx, caching, query helper,& test accounts
	s.setupGeneral()

}

func (s *KeeperTestHelper) setupGeneral() {
	s.setupGeneralCustomChainId("bitsong-2b")
}

func (s *KeeperTestHelper) setupGeneralCustomChainId(chainId string) {
	s.Ctx = s.App.BaseApp.NewContextLegacy(false,
		cmtproto.Header{Height: 1, ChainID: chainId, Time: defaultTestStartTime})
	if s.withCaching {
		s.Ctx, _ = s.Ctx.CacheContext()
	}
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = []sdk.AccAddress{}
	s.TestAccs = append(s.TestAccs, baseTestAccts...)
	s.hasUsedAbci = false
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

type GenerateAccountStrategy func(int) []sdk.AccAddress
type BondDenomProvider interface {
	BondDenom(ctx sdk.Context) string
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(bankKeeper bankkeeper.Keeper, stakingKeeper BondDenomProvider, ctx sdk.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
	return addTestAddrs(bankKeeper, stakingKeeper, ctx, accNum, accAmt, CreateRandomAccounts)
}

// addTestAddrs adds an account to the tests.
func addTestAddrs(bankKeeper bankkeeper.Keeper, stakingKeeper BondDenomProvider, ctx sdk.Context, accNum int, accAmt math.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)
	initCoins := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(bankKeeper, ctx, addr, initCoins)
	}

	return testAddrs
}

func initAccountWithCoins(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, coins); err != nil {
		panic(err)
	}

	if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins); err != nil {
		panic(err)
	}
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

// createTestConfigFiles creates minimal config files needed for upgrade handler tests
func (s *KeeperTestHelper) createTestConfigFiles(homeDir string) error {
	configDir := filepath.Join(homeDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}
	configContent := `
# Minimal config.toml for testing
[consensus]
timeout_commit = "5s"
timeout_propose = "3s"
timeout_propose_delta = "500ms"
timeout_prevote = "1s"
timeout_prevote_delta = "500ms"
timeout_precommit = "1s"
timeout_precommit_delta = "500ms"
 
`

	configPath := filepath.Join(configDir, "config.toml")
	return os.WriteFile(configPath, []byte(configContent), 0o644)
}
