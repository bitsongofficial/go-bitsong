package testing

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/v018/app"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestHelper struct {
	suite.Suite

	App         *app.BitsongApp
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper

	TestAccs []sdk.AccAddress
}

func (s *KeeperTestHelper) Setup() {
	s.App = app.Setup(false)
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "bitsong-test-suite-1", Time: time.Now().UTC()})
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = CreateRandomAccounts(3)
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
