package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

const (
	isCheckTx = false
)

var (
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	uri      = "ipfs://"
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initAmt)}
	symbol   = "btc"
	name     = "Bitcoin Network"

	maxSupply = sdk.NewInt(200000000)
	mintable  = true
	height    = int64(1)
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	bk          bankkeeper.Keeper
	app         *simapp.BitsongApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.FanTokenKeeper
	suite.bk = app.BankKeeper
	suite.app = app

	// set params
	suite.keeper.SetParamSet(suite.ctx, tokentypes.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, tokentypes.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, tokentypes.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setFanToken(token tokentypes.FanToken) {
	err := suite.keeper.AddFanToken(suite.ctx, token)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) issueFanToken(token tokentypes.FanToken) {
	suite.setFanToken(token)
}

func (suite *KeeperTestSuite) TestIssueFanToken() {
	denom, err := suite.keeper.IssueFanToken(suite.ctx, name, symbol, uri, maxSupply, owner)
	suite.NoError(err)

	suite.True(suite.keeper.HasFanToken(suite.ctx, denom))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)

	suite.Equal(owner, issuedToken.GetOwner())
	suite.Equal(uri, issuedToken.GetUri())
	// TODO: add more fields
}

func (suite *KeeperTestSuite) TestEditFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.setFanToken(fantokenObj)

	err := suite.keeper.EditFanToken(suite.ctx, fantokenObj.GetDenom(), false, owner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, fantokenObj.GetDenom())
	suite.NoError(err)

	fantokenObj.Mintable = false
	fantokenObj.MaxSupply = sdk.ZeroInt()
	suite.EqualValues(newToken.(*tokentypes.FanToken), &fantokenObj)
}

func (suite *KeeperTestSuite) TestMintFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	amt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("0ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, fantokenObj.GetDenom(), mintAmount, fantokenObj.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("1000ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())

	// mint token without owner

	err = suite.keeper.MintFanToken(suite.ctx, owner, fantokenObj.GetDenom(), mintAmount, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")
}

func (suite *KeeperTestSuite) TestBurnFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.issueFanToken(fantokenObj)

	amt := suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("0ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, fantokenObj.GetDenom(), mintAmount, fantokenObj.GetOwner())
	suite.NoError(err)

	burnedAmount := sdk.NewInt(200)

	err = suite.keeper.BurnFanToken(suite.ctx, fantokenObj.GetDenom(), burnedAmount, fantokenObj.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, fantokenObj.GetOwner(), fantokenObj.GetDenom())
	suite.Equal("800ft67188403047108DE19A31BEF25C8DABC1B6DC39B", amt.String())
}

func (suite *KeeperTestSuite) TestTransferFanToken() {
	fantokenObj := tokentypes.NewFanToken(name, symbol, uri, maxSupply, owner, height)
	suite.setFanToken(fantokenObj)

	dstOwner := sdk.AccAddress(tmhash.SumTruncated([]byte("TokenDstOwner")))

	err := suite.keeper.TransferFanTokenOwner(suite.ctx, fantokenObj.GetDenom(), fantokenObj.GetOwner(), dstOwner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, fantokenObj.GetDenom())
	suite.NoError(err)

	suite.Equal(dstOwner, newToken.GetOwner())
}
