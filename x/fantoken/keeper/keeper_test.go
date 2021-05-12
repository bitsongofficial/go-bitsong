package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	simapp "github.com/bitsongofficial/bitsong/app"
	"github.com/bitsongofficial/bitsong/types"
	"github.com/bitsongofficial/bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/bitsong/x/fantoken/types"
)

const (
	isCheckTx = false
)

var (
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(types.BondDenom, initAmt)}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	bk          bankkeeper.Keeper
	app         *simapp.Bitsong
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
	token := tokentypes.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), false, "test", owner)

	err := suite.keeper.IssueFanToken(
		suite.ctx, token.Symbol, token.Name,
		token.MaxSupply, token.Mintable, token.Description, token.GetOwner(),
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasFanToken(suite.ctx, token.Symbol))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwner().String())

	suite.EqualValues(&token, issuedToken.(*tokentypes.FanToken))
}

func (suite *KeeperTestSuite) TestUpdateFanTokenMintable() {
	token := tokentypes.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), false, "test", owner)
	suite.setFanToken(token)

	symbol := "btc"
	mintable := true

	err := suite.keeper.UpdateFanTokenMintable(suite.ctx, symbol, mintable, owner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, symbol)
	suite.NoError(err)

	expToken := tokentypes.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), mintable, "test", owner)

	suite.EqualValues(newToken.(*tokentypes.FanToken), &expToken)
}

func (suite *KeeperTestSuite) TestMintFanToken() {
	token := tokentypes.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(2000), true, "test", owner)
	suite.issueFanToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Symbol)
	suite.Equal("0btc", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, token.Symbol, mintAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("1000000000ubtc", amt.String())

	// mint token without owner

	err = suite.keeper.MintFanToken(suite.ctx, owner, token.Symbol, mintAmount, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")

	token = tokentypes.NewFanToken("atom", "Cosmos Hub", sdk.NewInt(2000), true, "test", sdk.AccAddress{})
	suite.issueFanToken(token)

	err = suite.keeper.MintFanToken(suite.ctx, owner, token.Symbol, mintAmount, sdk.AccAddress{})
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestBurnFanToken() {
	token := tokentypes.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(2000), true, "test", owner)
	suite.issueFanToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Symbol)
	suite.Equal("0btc", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, token.Symbol, mintAmount, token.GetOwner())
	suite.NoError(err)

	burnedAmount := sdk.NewInt(200)

	err = suite.keeper.BurnFanToken(suite.ctx, token.Symbol, burnedAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetDenom())
	suite.Equal("800000000ubtc", amt.String())
}

func (suite *KeeperTestSuite) TestTransferFanToken() {
	token := tokentypes.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), false, "test", owner)
	suite.setFanToken(token)

	dstOwner := sdk.AccAddress(tmhash.SumTruncated([]byte("TokenDstOwner")))

	err := suite.keeper.TransferFanTokenOwner(suite.ctx, token.Symbol, token.GetOwner(), dstOwner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.Equal(dstOwner, newToken.GetOwner())
}
