package authenticator_test

import (
	"os"
	"testing"
	"time"

	"cosmossdk.io/store/prefix"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/smart-account/ante"
	"github.com/bitsongofficial/go-bitsong/x/smart-account/authenticator"
	"github.com/bitsongofficial/go-bitsong/x/smart-account/post"
	"github.com/bitsongofficial/go-bitsong/x/smart-account/testutils"

	storetypes "cosmossdk.io/store/types"
)

type SpendLimitAuthenticatorTest struct {
	BaseAuthenticatorSuite

	Store                      prefix.Store
	CosmwasmAuth               authenticator.CosmwasmAuthenticator
	AlwaysPassAuth             testutils.TestingAuthenticator
	AuthenticatorAnteDecorator ante.AuthenticatorDecorator
	AuthenticatorPostDecorator post.AuthenticatorPostDecorator
}

type InstantiateMsg struct {
	PriceResolutionConfig PriceResolutionConfig `json:"price_resolution_config"`
	TrackedDenoms         []TrackedDenom        `json:"tracked_denoms"`
}

type TrackedDenom struct {
	Denom      string              `json:"denom"`
	SwapRoutes []SwapAmountInRoute `json:"swap_routes"`
}

type SwapAmountInRoute struct {
	PoolID        string `json:"pool_id"` // as u64
	TokenOutDenom string `json:"token_out_denom"`
}

type PriceResolutionConfig struct {
	QuoteDenom         string `json:"quote_denom"`
	StalenessThreshold string `json:"staleness_threshold"` // as u64
	TwapDuration       string `json:"twap_duration"`       // as u64
}

// params
type SpendLimitParams struct {
	Limit       string     `json:"limit"`        // as u128
	ResetPeriod string     `json:"reset_period"` // day | week | month | year
	TimeLimit   *TimeLimit `json:"time_limit,omitempty"`
}

type TimeLimit struct {
	Start *string `json:"start,omitempty"` // as u64 or None
	End   string  `json:"end"`             // as u64
}

func TestSpendLimitAuthenticatorTest(t *testing.T) {
	suite.Run(t, new(SpendLimitAuthenticatorTest))
}

const UUSDC = "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4"

func (s *SpendLimitAuthenticatorTest) SetupTest() {
	s.SetupKeys()
	s.BitsongApp = app.Setup(false)
	s.Ctx = s.BitsongApp.NewContextLegacy(false, tmproto.Header{})
	s.Ctx = s.Ctx.WithGasMeter(storetypes.NewGasMeter(10_000_000))
	s.Ctx = s.Ctx.WithBlockTime(time.Now())
	s.EncodingConfig = app.MakeEncodingConfig()

	s.CosmwasmAuth = authenticator.NewCosmwasmAuthenticator(s.BitsongApp.ContractKeeper, s.BitsongApp.AccountKeeper, s.BitsongApp.AppCodec())

	s.AlwaysPassAuth = testutils.TestingAuthenticator{Approve: testutils.Always, Confirm: testutils.Always, GasConsumption: 0}
	s.BitsongApp.SmartAccountKeeper.AuthenticatorManager.RegisterAuthenticator(s.AlwaysPassAuth)

	deductFeeDecorator := sdkante.NewDeductFeeDecorator(s.BitsongApp.AccountKeeper, s.BitsongApp.BankKeeper, s.BitsongApp.FeeGrantKeeper, nil)
	s.AuthenticatorAnteDecorator = ante.NewAuthenticatorDecorator(
		s.BitsongApp.AppCodec(),
		s.BitsongApp.SmartAccountKeeper,
		s.BitsongApp.AccountKeeper,
		s.EncodingConfig.TxConfig.SignModeHandler(),
		deductFeeDecorator,
	)

	s.AuthenticatorPostDecorator = post.NewAuthenticatorPostDecorator(
		s.BitsongApp.AppCodec(),
		s.BitsongApp.SmartAccountKeeper,
		s.BitsongApp.AccountKeeper,
		s.EncodingConfig.TxConfig.SignModeHandler(),
		// Add an empty handler here to enable a circuit breaker pattern
		sdk.ChainPostDecorators(sdk.Terminator{}), //nolint
	)
}

func (s *SpendLimitAuthenticatorTest) TearDownTest() {
	os.RemoveAll(s.HomeDir)
}

func (s *SpendLimitAuthenticatorTest) TestSpendLimit() {
	// anteHandler := sdk.ChainAnteDecorators(s.AuthenticatorAnteDecorator)
	// postHandler := sdk.ChainPostDecorators(s.AuthenticatorPostDecorator)

	// increase time by 1hr to ensure twap price is available
	s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(time.Hour))

	// bz, err := json.Marshal(msg)
	// s.Require().NoError(err)
	// contractAddr := s.InstantiateContract(string(bz), codeId)

	// add new authenticator
	// ak := s.BitsongApp.SmartAccountKeeper

	// authAcc := s.TestAccAddress[1]
	// authAccPriv := s.TestPrivKeys[1]

	// initData := authenticator.CosmwasmAuthenticatorInitData{
	// 	Contract: contractAddr.String(),
	// 	Params:   bz,
	// }

	// bz, err = json.Marshal(initData)
	// s.Require().NoError(err)

	// // hack to get fee payer authenticated
	// id, err := ak.AddAuthenticator(s.Ctx, authAcc, s.AlwaysPassAuth.Type(), []byte{})
	// s.Require().NoError(err)
	// s.Require().Equal(uint64(1), id)

	// id, err = ak.AddAuthenticator(s.Ctx, authAcc, authenticator.CosmwasmAuthenticator{}.Type(), bz)
	// s.Require().NoError(err)
	// s.Require().Equal(uint64(2), id)

	// // fund acc
	// s.FundAcc(authAcc, sdk.NewCoins(sdk.NewCoin(UUSDC, osmomath.NewInt(200000000000))))
	// s.FundAcc(authAcc, sdk.NewCoins(sdk.NewCoin(appparams.CoinUnit, osmomath.NewInt(200000000000))))

	// // a hack for setting fee payer
	// selfSend := banktypes.MsgSend{
	// 	FromAddress: authAcc.String(),
	// 	ToAddress:   authAcc.String(),
	// 	Amount:      sdk.NewCoins(sdk.NewCoin(UUSDC, osmomath.NewInt(1))),
	// }

	//... somthing here

	// ante
	// _, err = anteHandler(s.Ctx, tx, false)
	// s.Require().NoError(err)

	//... somthing here

	// post
	// _, err = postHandler(s.Ctx, tx, false, true)
	// s.Require().NoError(err)

	//... somthing here

	// ante
	// _, err = anteHandler(s.Ctx, tx, false)
	// s.Require().NoError(err)

	//... somthing here
	// swap

	// post
	// _, err = postHandler(s.Ctx, tx, false, true)
	// s.Require().Error(err)

	//... fix error

	// ante
	// _, err = anteHandler(s.Ctx, tx, false)
	// s.Require().NoError(err)

	//... somthing here

	// post
	// _, err = postHandler(s.Ctx, tx, false, true)
	// s.Require().NoError(err)

	//... somthing here

	// ante
	// _, err = anteHandler(s.Ctx, tx, false)
	// s.Require().Error(err)

}

func (s *SpendLimitAuthenticatorTest) StoreContractCode(path string) uint64 {
	btsgApp := s.BitsongApp
	govKeeper := wasmkeeper.NewGovPermissionKeeper(btsgApp.WasmKeeper)
	creator := btsgApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	wasmCode, err := os.ReadFile(path)
	s.Require().NoError(err)
	accessEveryone := wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody}
	codeID, _, err := govKeeper.Create(s.Ctx, creator, wasmCode, &accessEveryone)
	s.Require().NoError(err)
	return codeID
}

func (s *SpendLimitAuthenticatorTest) InstantiateContract(msg string, codeID uint64) sdk.AccAddress {
	btsgApp := s.BitsongApp
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(btsgApp.WasmKeeper)
	creator := btsgApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	addr, _, err := contractKeeper.Instantiate(s.Ctx, codeID, creator, creator, []byte(msg), "contract", nil)
	s.Require().NoError(err)
	return addr
}
