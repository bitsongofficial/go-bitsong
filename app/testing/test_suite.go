package apptesting

import (
	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakinghelper "github.com/cosmos/cosmos-sdk/x/staking/testutil"

	"github.com/stretchr/testify/suite"
)

type KeeperTestHelper struct {
	suite.Suite

	App         *app.BitsongApp                 // Mock bitsong application
	Ctx         sdk.Context                     // simulated context
	QueryHelper *baseapp.QueryServiceTestHelper // GRPC query simulator
	TestAccs    []sdk.AccAddress                // Test accounts

	StakingHelper *stakinghelper.Helper // Useful staking helpers
}

func (s *KeeperTestHelper) Setup() {
	s.App = app.Setup(s.T())
	s.Ctx = s.App.BaseApp.NewContext(false)
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = CreateRandomAccounts(3)
	s.StakingHelper = stakinghelper.NewHelper(s.Suite.T(), s.Ctx, s.App.AppKeepers.StakingKeeper)
	s.StakingHelper.Denom = "ubtsg"
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

// FundAcc funds target address with specified amount.
func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := testutil.FundAccount(s.Ctx, s.App.AppKeepers.BankKeeper, acc, amounts)
	s.Require().NoError(err)
}
