package keeper_test

import (
	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/suite"
)

var (
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	uri      = "ipfs://"
	initAmt  = sdk.NewIntFromUint64(100)
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
	app := simapp.Setup(false)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = app.FanTokenKeeper
	suite.bk = app.BankKeeper
	suite.app = app

	// set params
	suite.keeper.SetParamSet(suite.ctx, fantokentypes.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, fantokentypes.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, fantokentypes.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestIssue() {
	denom, err := suite.keeper.Issue(suite.ctx, name, symbol, uri, maxSupply, owner, owner)
	suite.NoError(err)
	suite.True(suite.keeper.HasFanToken(suite.ctx, denom))

	issuedToken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)

	suite.Equal(denom, issuedToken.GetDenom())
	suite.Equal(name, issuedToken.GetName())
	suite.Equal(symbol, issuedToken.GetSymbol())
	suite.Equal(uri, issuedToken.GetURI())
	suite.Equal(maxSupply, issuedToken.GetMaxSupply())
	suite.Equal(owner, issuedToken.GetAuthority())
	suite.Equal(owner, issuedToken.GetMinter())
}

func (suite *KeeperTestSuite) TestMint() {
	// issue a new fantoken
	denom, err := suite.keeper.Issue(suite.ctx, name, symbol, uri, maxSupply, owner, owner)
	suite.NoError(err)

	// check actual fantoken balance
	balance := suite.bk.GetBalance(suite.ctx, owner, denom)
	suite.Equal(sdk.ZeroInt(), balance.Amount)
	suite.NotEqual(sdk.NewInt(2), balance.Amount)

	// check the fantoken supply
	supply := suite.bk.GetSupply(suite.ctx, denom)
	suite.Equal(sdk.ZeroInt(), supply.Amount)
	suite.Equal(denom, supply.Denom)

	// mint some token
	suite.keeper.Mint(suite.ctx, owner, owner, sdk.NewCoin(denom, sdk.NewInt(10)))

	// check the fantoken balance once a time
	balance = suite.bk.GetBalance(suite.ctx, owner, denom)
	suite.Equal(sdk.NewInt(10), balance.Amount)
	suite.NotEqual(sdk.ZeroInt(), balance.Amount)

	// check the fantoken supply once a time
	supply = suite.bk.GetSupply(suite.ctx, denom)
	suite.Equal(sdk.NewInt(10), supply.Amount)
	suite.Equal(denom, supply.Denom)
}

func (suite *KeeperTestSuite) TestBurn() {
	// issue a new fantoken
	denom, err := suite.keeper.Issue(suite.ctx, name, symbol, uri, maxSupply, owner, owner)
	suite.NoError(err)

	// mint some token
	suite.keeper.Mint(suite.ctx, owner, owner, sdk.NewCoin(denom, sdk.NewInt(10)))

	// burn some token
	suite.keeper.Burn(suite.ctx, sdk.NewCoin(denom, sdk.NewInt(6)), owner)

	// check the fantoken balance
	balance := suite.bk.GetBalance(suite.ctx, owner, denom)
	suite.Equal(sdk.NewInt(4), balance.Amount)
	suite.NotEqual(sdk.ZeroInt(), balance.Amount)

	// check the fantoken supply once a time
	supply := suite.bk.GetSupply(suite.ctx, denom)
	suite.Equal(sdk.NewInt(4), supply.Amount)
	suite.Equal(denom, supply.Denom)
}

func (suite *KeeperTestSuite) TestSetMinter() {
	// issue a new fantoken
	denom, err := suite.keeper.Issue(suite.ctx, name, symbol, uri, maxSupply, owner, owner)
	suite.NoError(err)

	// set the new minter
	err = suite.keeper.SetMinter(suite.ctx, denom, owner, owner)
	suite.NoError(err)

	// set an empty oldMinter
	err = suite.keeper.SetMinter(suite.ctx, denom, sdk.AccAddress{}, sdk.AccAddress{})
	suite.Error(err)

	// set an empty minter
	err = suite.keeper.SetMinter(suite.ctx, denom, owner, sdk.AccAddress{})
	suite.NoError(err)

	// after an empty minter, you cannot change the minter again
	err = suite.keeper.SetMinter(suite.ctx, denom, owner, sdk.AccAddress{})
	suite.Error(err)

	err = suite.keeper.SetMinter(suite.ctx, denom, sdk.AccAddress{}, sdk.AccAddress{})
	suite.Error(err)
}

func (suite *KeeperTestSuite) TestSetAuthority() {
	// issue a new fantoken
	denom, err := suite.keeper.Issue(suite.ctx, name, symbol, uri, maxSupply, owner, owner)
	suite.NoError(err)

	// set the new authority
	err = suite.keeper.SetAuthority(suite.ctx, denom, owner, owner)
	suite.NoError(err)

	// set an empty oldAuthority
	err = suite.keeper.SetAuthority(suite.ctx, denom, sdk.AccAddress{}, sdk.AccAddress{})
	suite.Error(err)

	// set an empty authority
	err = suite.keeper.SetAuthority(suite.ctx, denom, owner, sdk.AccAddress{})
	suite.NoError(err)

	// after an empty authority, you cannot change the authority again
	err = suite.keeper.SetAuthority(suite.ctx, denom, owner, sdk.AccAddress{})
	suite.Error(err)

	err = suite.keeper.SetAuthority(suite.ctx, denom, sdk.AccAddress{}, sdk.AccAddress{})
	suite.Error(err)
}

func (suite *KeeperTestSuite) TestSetUri() {
	// issue a new fantoken
	denom, err := suite.keeper.Issue(suite.ctx, name, symbol, uri, maxSupply, owner, owner)
	suite.NoError(err)

	newUri := "ipfs://newUri"

	// set the new uri
	err = suite.keeper.SetUri(suite.ctx, denom, newUri, owner)
	suite.NoError(err)

	// get fantoken
	fantoken, err := suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)
	suite.Equal(newUri, fantoken.GetURI())

	emptyUri := ""
	// set the new uri
	err = suite.keeper.SetUri(suite.ctx, denom, emptyUri, owner)
	suite.NoError(err)

	// check fantoken again
	fantoken, err = suite.keeper.GetFanToken(suite.ctx, denom)
	suite.NoError(err)
	suite.Equal(emptyUri, fantoken.GetURI())

	malformedUri := "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
		"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
		"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
		"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
		"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" +
		"0123456789012"

	// set the new uri
	err = suite.keeper.SetUri(suite.ctx, denom, malformedUri, owner)
	suite.Error(err)
}
