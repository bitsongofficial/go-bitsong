package keeper_test

import (
	"crypto/sha256"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/suite"

	_ "embed"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	"github.com/bitsongofficial/go-bitsong/x/cadance/keeper"
	"github.com/bitsongofficial/go-bitsong/x/cadance/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type IntegrationTestSuite struct {
	apptesting.KeeperTestHelper
	bk               minttypes.BankKeeper
	wk               wasmkeeper.Keeper
	cadanceMsgServer types.MsgServer
	queryClient      types.QueryClient
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	s.Setup()
	s.bk = s.App.AppKeepers.BankKeeper
	s.wk = s.App.AppKeepers.WasmKeeper

	// encCfg := app.MakeEncodingConfig()
	// types.RegisterInterfaces(encCfg.InterfaceRegistry)
	s.queryClient = types.NewQueryClient(s.QueryHelper)
	s.cadanceMsgServer = keeper.NewMsgServerImpl(s.App.AppKeepers.CadanceKeeper)
}

func (s *IntegrationTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := s.bk.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.bk.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

//go:embed testdata/clock_example.wasm
var wasmContract []byte

// stores  uploads wasm code in example
func (s *IntegrationTestSuite) StoreCode() {
	_, _, sender := testdata.KeyTestPubAddr()

	params := s.App.AppKeepers.WasmKeeper.GetParams(s.Ctx)
	params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
	err := s.App.AppKeepers.WasmKeeper.SetParams(s.Ctx, params)
	s.Require().NoError(err)
	msg := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender.String()
	})
	rsp, err := s.App.MsgServiceRouter().Handler(msg)(s.Ctx, msg)
	s.Require().NoError(err)
	var result wasmtypes.MsgStoreCodeResponse
	s.Require().NoError(s.App.AppCodec().Unmarshal(rsp.Data, &result))
	s.Require().Equal(uint64(1), result.CodeID)
	expHash := sha256.Sum256(wasmContract)
	s.Require().Equal(expHash[:], result.Checksum)
	// and
	info := s.App.AppKeepers.WasmKeeper.GetCodeInfo(s.Ctx, 1)
	s.Require().NotNil(info)
	s.Require().Equal(expHash[:], info.CodeHash)
	s.Require().Equal(sender.String(), info.Creator)
	s.Require().Equal(wasmtypes.DefaultParams().InstantiateDefaultPermission.With(sender), info.InstantiateConfig)
}

func (s *IntegrationTestSuite) InstantiateContract(sender string, admin string) string {
	msgStoreCode := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender
	})
	_, err := s.App.MsgServiceRouter().Handler(msgStoreCode)(s.Ctx, msgStoreCode)
	s.Require().NoError(err)

	msgInstantiate := wasmtypes.MsgInstantiateContractFixture(func(m *wasmtypes.MsgInstantiateContract) {
		m.Sender = sender
		m.Admin = admin
		m.Msg = []byte(`{}`)
	})
	resp, err := s.App.MsgServiceRouter().Handler(msgInstantiate)(s.Ctx, msgInstantiate)
	s.Require().NoError(err)
	var result wasmtypes.MsgInstantiateContractResponse
	s.Require().NoError(s.App.AppCodec().Unmarshal(resp.Data, &result))
	contractInfo := s.App.AppKeepers.WasmKeeper.GetContractInfo(s.Ctx, sdk.MustAccAddressFromBech32(result.Address))
	s.Require().Equal(contractInfo.CodeID, uint64(1))
	s.Require().Equal(contractInfo.Admin, admin)
	s.Require().Equal(contractInfo.Creator, sender)

	return result.Address
}

// Helper method for quickly registering a cadance contract
func (s *IntegrationTestSuite) RegisterCadanceContract(senderAddress string, contractAddress string) {
	err := s.App.AppKeepers.CadanceKeeper.RegisterContract(s.Ctx, senderAddress, contractAddress)
	s.Require().NoError(err)
}

// Helper method for quickly unregistering a cadance contract
func (s *IntegrationTestSuite) UnregisterCadanceContract(senderAddress string, contractAddress string) {
	err := s.App.AppKeepers.CadanceKeeper.UnregisterContract(s.Ctx, senderAddress, contractAddress)
	s.Require().NoError(err)
}

// Helper method for quickly jailing a cadance contract
func (s *IntegrationTestSuite) JailCadanceContract(contractAddress string) {
	err := s.App.AppKeepers.CadanceKeeper.SetJailStatus(s.Ctx, contractAddress, true)
	s.Require().NoError(err)
}

// Helper method for quickly unjailing a cadance contract
func (s *IntegrationTestSuite) UnjailCadanceContract(senderAddress string, contractAddress string) {
	err := s.App.AppKeepers.CadanceKeeper.SetJailStatusBySender(s.Ctx, senderAddress, contractAddress, false)
	s.Require().NoError(err)
}
