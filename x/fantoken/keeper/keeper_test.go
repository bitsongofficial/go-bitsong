package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/bitsongofficial/ledger/simapp"
	"github.com/bitsongofficial/ledger/x/fantoken/keeper"
	"github.com/bitsongofficial/ledger/x/fantoken/types"
)

const (
	isCheckTx = false
)

var (
	denom    = types.GetNativeToken().Denom
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	bk          bankkeeper.Keeper
	app         *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper
	suite.app = app

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setFanToken(token types.FanToken) {
	err := suite.keeper.AddFanToken(suite.ctx, token)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) issueFanToken(token types.FanToken) {
	suite.setFanToken(token)
}

func (suite *KeeperTestSuite) TestIssueFanToken() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), false, owner)

	err := suite.keeper.IssueFanToken(
		suite.ctx, token.Denom, token.Name,
		token.MaxSupply, token.Mintable, token.GetOwner(),
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasFanToken(suite.ctx, token.Denom))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, token.Denom)
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwner().String())

	suite.EqualValues(&token, issuedToken.(*types.FanToken))
}

func (suite *KeeperTestSuite) TestUpdateFanTokenMintable() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), false, owner)
	suite.setFanToken(token)

	denom := "btc"
	mintable := true

	err := suite.keeper.UpdateFanTokenMintable(suite.ctx, denom, mintable, owner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)

	expToken := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), mintable, owner)

	suite.EqualValues(newToken.(*types.FanToken), &expToken)
}

func (suite *KeeperTestSuite) TestMintFanToken() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(2000), true, owner)
	suite.issueFanToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom)
	suite.Equal("0btc", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, token.Denom, mintAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom)
	suite.Equal("1000btc", amt.String())

	// mint token without owner

	err = suite.keeper.MintFanToken(suite.ctx, owner, token.Denom, mintAmount, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")

	token = types.NewFanToken("atom", "Cosmos Hub", sdk.NewInt(2000), true, sdk.AccAddress{})
	suite.issueFanToken(token)

	err = suite.keeper.MintFanToken(suite.ctx, owner, token.Denom, mintAmount, sdk.AccAddress{})
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestBurnFanToken() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(2000), true, owner)
	suite.issueFanToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom)
	suite.Equal("0btc", amt.String())

	mintAmount := sdk.NewInt(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintFanToken(suite.ctx, recipient, token.Denom, mintAmount, token.GetOwner())
	suite.NoError(err)

	burnedAmount := sdk.NewInt(200)

	err = suite.keeper.BurnFanToken(suite.ctx, token.Denom, burnedAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.Denom)
	suite.Equal("800btc", amt.String())
}

func (suite *KeeperTestSuite) TestTransferFanToken() {
	token := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(21000000), false, owner)
	suite.setFanToken(token)

	dstOwner := sdk.AccAddress(tmhash.SumTruncated([]byte("TokenDstOwner")))

	err := suite.keeper.TransferFanTokenOwner(suite.ctx, token.Denom, token.GetOwner(), dstOwner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetFanToken(suite.ctx, token.Denom)
	suite.NoError(err)

	suite.Equal(dstOwner, newToken.GetOwner())
}
