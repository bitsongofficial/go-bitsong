package cadance_test

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	"github.com/stretchr/testify/suite"

	_ "embed"

	cadance "github.com/bitsongofficial/go-bitsong/x/cadance"
	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type EndBlockerTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestEndBlockerTestSuite(t *testing.T) {
	suite.Run(t, new(EndBlockerTestSuite))
}

func (s *EndBlockerTestSuite) SetupTest() {
	s.Setup()

}

//go:embed keeper/testdata/clock_example.wasm
var cadanceContract []byte

//go:embed keeper/testdata/cw_testburn.wasm
var burnContract []byte

func (s *EndBlockerTestSuite) StoreCode(wasmContract []byte) {
	_, _, sender := testdata.KeyTestPubAddr()
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
	info := s.App.WasmKeeper.GetCodeInfo(s.Ctx, 1)
	s.Require().NotNil(info)
	s.Require().Equal(expHash[:], info.CodeHash)
	s.Require().Equal(sender.String(), info.Creator)
	s.Require().Equal(wasmtypes.DefaultParams().InstantiateDefaultPermission.With(sender), info.InstantiateConfig)
}

func (s *EndBlockerTestSuite) InstantiateContract(sender string, admin string) string {
	msgStoreCode := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = cadanceContract
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
	contractInfo := s.App.WasmKeeper.GetContractInfo(s.Ctx, sdk.MustAccAddressFromBech32(result.Address))
	s.Require().Equal(contractInfo.CodeID, uint64(1))
	s.Require().Equal(contractInfo.Admin, admin)
	s.Require().Equal(contractInfo.Creator, sender)

	return result.Address
}

func (s *EndBlockerTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := s.App.BankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

// Register a contract. You must store the contract code before registering.
func (s *EndBlockerTestSuite) registerContract() string {
	// Create & fund accounts
	_, _, sender := testdata.KeyTestPubAddr()
	_, _, admin := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, sender, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))
	_ = s.FundAccount(s.Ctx, admin, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	// Instantiate contract
	contractAddress := s.InstantiateContract(sender.String(), admin.String())

	// Register contract
	cadanceKeeper := s.App.CadanceKeeper
	err := cadanceKeeper.RegisterContract(s.Ctx, admin.String(), contractAddress)
	s.Require().NoError(err)

	// Assert contract is registered
	contract, err := cadanceKeeper.GetCadanceContract(s.Ctx, contractAddress)
	s.Require().NoError(err)
	s.Require().Equal(contractAddress, contract.ContractAddress)

	// Increment block height
	s.Ctx = s.Ctx.WithBlockHeight(11)

	return contract.ContractAddress
}

// Test the end blocker. This test registers a contract, executes it with enough gas,
// too little gas, and also ensures the unjailing process functions.
func (s *EndBlockerTestSuite) TestEndBlocker() {
	// Setup test
	cadanceKeeper := s.App.CadanceKeeper
	s.StoreCode(cadanceContract)
	contractAddress := s.registerContract()

	// Query contract
	val := s.queryContract(contractAddress)
	s.Require().Equal(int64(0), val)

	// Call end blocker
	s.callEndBlocker()

	// Query contract
	val = s.queryContract(contractAddress)
	s.Require().Equal(int64(1), val)

	// Update params with 10 gas limit
	s.updateGasLimit(65_000)

	// Call end blocker
	s.callEndBlocker()

	// Ensure contract is now jailed
	contract, err := cadanceKeeper.GetCadanceContract(s.Ctx, contractAddress)
	s.Require().NoError(err)
	s.Require().True(contract.IsJailed)

	// Update params to regular
	s.updateGasLimit(types.DefaultParams().ContractGasLimit)

	// Call end blocker
	s.callEndBlocker()

	// Unjail contract
	err = cadanceKeeper.SetJailStatus(s.Ctx, contractAddress, false)
	s.Require().NoError(err)

	// Ensure contract is no longer jailed
	contract, err = cadanceKeeper.GetCadanceContract(s.Ctx, contractAddress)
	s.Require().NoError(err)
	s.Require().False(contract.IsJailed)

	// Call end blocker
	s.callEndBlocker()

	// Query contract
	val = s.queryContract(contractAddress)
	s.Require().Equal(int64(2), val)
}

// Test a contract which does not handle the sudo EndBlock msg.
func (s *EndBlockerTestSuite) TestInvalidContract() {
	// Setup test
	cadanceKeeper := s.App.CadanceKeeper
	s.StoreCode(burnContract)
	contractAddress := s.registerContract()

	// Run the end blocker
	s.callEndBlocker()

	// Ensure contract is now jailed
	contract, err := cadanceKeeper.GetCadanceContract(s.Ctx, contractAddress)
	s.Require().NoError(err)
	s.Require().True(contract.IsJailed)
}

// Test the endblocker with numerous contracts that all panic
func (s *EndBlockerTestSuite) TestPerformance() {
	s.StoreCode(burnContract)

	numContracts := 1000

	// Register numerous contracts
	for x := 0; x < numContracts; x++ {
		// Register contract
		_ = s.registerContract()
	}

	// Ensure contracts exist
	cadanceKeeper := s.App.CadanceKeeper
	contracts, err := cadanceKeeper.GetAllContracts(s.Ctx)
	s.Require().NoError(err)
	s.Require().Len(contracts, numContracts)

	// Call end blocker
	s.callEndBlocker()

	// Ensure contracts are jailed
	contracts, err = cadanceKeeper.GetAllContracts(s.Ctx)
	s.Require().NoError(err)
	for _, contract := range contracts {
		s.Require().True(contract.IsJailed)
	}
}

// Update the gas limit
func (s *EndBlockerTestSuite) updateGasLimit(gasLimit uint64) {
	params := types.DefaultParams()
	params.ContractGasLimit = gasLimit
	k := s.App.CadanceKeeper

	store := s.Ctx.KVStore(k.GetStore())
	bz := k.GetCdc().MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
}

// Call the end blocker, incrementing the block height
func (s *EndBlockerTestSuite) callEndBlocker() {
	cadance.EndBlocker(s.Ctx, s.App.CadanceKeeper)
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
}

// Query the cadance contract
func (s *EndBlockerTestSuite) queryContract(contractAddress string) int64 {
	query := `{"get_config":{}}`
	output, err := s.App.WasmKeeper.QuerySmart(s.Ctx, sdk.MustAccAddressFromBech32(contractAddress), []byte(query))
	s.Require().NoError(err)

	var val struct {
		Val int64 `json:"val"`
	}

	err = json.Unmarshal(output, &val)
	s.Require().NoError(err)

	return val.Val
}
