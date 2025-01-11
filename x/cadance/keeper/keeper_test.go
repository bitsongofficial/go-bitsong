package keeper_test

import (
	"crypto/sha256"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/suite"

	_ "embed"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/cadance/keeper"
	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	app              *app.BitsongApp
	bankKeeper       bankkeeper.Keeper
	queryClient      types.QueryClient
	cadanceMsgServer types.MsgServer
}

func (s *IntegrationTestSuite) SetupTest() {
	isCheckTx := false
	s.app = app.Setup(s.T())

	s.ctx = s.app.BaseApp.NewContext(isCheckTx)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(s.app.AppKeepers.CadanceKeeper))

	s.queryClient = types.NewQueryClient(queryHelper)
	s.bankKeeper = s.app.AppKeepers.BankKeeper
	s.cadanceMsgServer = keeper.NewMsgServerImpl(s.app.AppKeepers.CadanceKeeper)
}

func (s *IntegrationTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := s.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

//go:embed testdata/clock_example.wasm
var wasmContract []byte

// stores  uploads wasm code in example
func (s *IntegrationTestSuite) StoreCode() {
	_, _, sender := testdata.KeyTestPubAddr()
	msg := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender.String()
	})
	rsp, err := s.app.MsgServiceRouter().Handler(msg)(s.ctx, msg)
	s.Require().NoError(err)
	var result wasmtypes.MsgStoreCodeResponse
	s.Require().NoError(s.app.AppCodec().Unmarshal(rsp.Data, &result))
	s.Require().Equal(uint64(1), result.CodeID)
	expHash := sha256.Sum256(wasmContract)
	s.Require().Equal(expHash[:], result.Checksum)
	// and
	info := s.app.AppKeepers.WasmKeeper.GetCodeInfo(s.ctx, 1)
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
	_, err := s.app.MsgServiceRouter().Handler(msgStoreCode)(s.ctx, msgStoreCode)
	s.Require().NoError(err)

	msgInstantiate := wasmtypes.MsgInstantiateContractFixture(func(m *wasmtypes.MsgInstantiateContract) {
		m.Sender = sender
		m.Admin = admin
		m.Msg = []byte(`{}`)
	})
	resp, err := s.app.MsgServiceRouter().Handler(msgInstantiate)(s.ctx, msgInstantiate)
	s.Require().NoError(err)
	var result wasmtypes.MsgInstantiateContractResponse
	s.Require().NoError(s.app.AppCodec().Unmarshal(resp.Data, &result))
	contractInfo := s.app.AppKeepers.WasmKeeper.GetContractInfo(s.ctx, sdk.MustAccAddressFromBech32(result.Address))
	s.Require().Equal(contractInfo.CodeID, uint64(1))
	s.Require().Equal(contractInfo.Admin, admin)
	s.Require().Equal(contractInfo.Creator, sender)

	return result.Address
}

// Helper method for quickly registering a cadance contract
func (s *IntegrationTestSuite) RegisterCadanceContract(senderAddress string, contractAddress string) {
	err := s.app.AppKeepers.CadanceKeeper.RegisterContract(s.ctx, senderAddress, contractAddress)
	s.Require().NoError(err)
}

// Helper method for quickly unregistering a cadance contract
func (s *IntegrationTestSuite) UnregisterCadanceContract(senderAddress string, contractAddress string) {
	err := s.app.AppKeepers.CadanceKeeper.UnregisterContract(s.ctx, senderAddress, contractAddress)
	s.Require().NoError(err)
}

// Helper method for quickly jailing a cadance contract
func (s *IntegrationTestSuite) JailCadanceContract(contractAddress string) {
	err := s.app.AppKeepers.CadanceKeeper.SetJailStatus(s.ctx, contractAddress, true)
	s.Require().NoError(err)
}

// Helper method for quickly unjailing a cadance contract
func (s *IntegrationTestSuite) UnjailCadanceContract(senderAddress string, contractAddress string) {
	err := s.app.AppKeepers.CadanceKeeper.SetJailStatusBySender(s.ctx, senderAddress, contractAddress, false)
	s.Require().NoError(err)
}
