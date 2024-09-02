package app_test

import (
	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
	"time"
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())

	valKey  = ed25519.GenPrivKey()
	valAddr = sdk.AccAddress(valKey.PubKey().Address())

	zeroCommissionRates = types.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	oneDec, _           = sdk.NewDecFromStr("0.01")
	tenDec, _           = sdk.NewDecFromStr("0.10")
	twentyDec, _        = sdk.NewDecFromStr("0.20")
	goodCommissionRates = types.NewCommissionRates(tenDec, twentyDec, oneDec)

	PKs = simapp.CreateTestPubKeys(500)
)

func TestCreateValidatorFailAnteHandler(t *testing.T) {
	genTokens := sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	bondTokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

	acc1 := &authtypes.BaseAccount{Address: addr1.String()}
	acc2 := &authtypes.BaseAccount{Address: addr2.String()}
	accs := authtypes.GenesisAccounts{acc1, acc2}
	balances := []banktypes.Balance{
		{
			Address: addr1.String(),
			Coins:   sdk.Coins{genCoin},
		},
		{
			Address: addr2.String(),
			Coins:   sdk.Coins{genCoin},
		},
	}

	app := simapp.SetupWithGenesisAccounts(accs, balances...)
	simapp.CheckBalance(t, app, addr1, sdk.Coins{genCoin})
	simapp.CheckBalance(t, app, addr2, sdk.Coins{genCoin})

	// create validator
	description := types.NewDescription("foo_moniker", "", "", "", "")
	createValidatorMsg, err := types.NewMsgCreateValidator(
		sdk.ValAddress(addr1), valKey.PubKey(), bondCoin, description, zeroCommissionRates, sdk.OneInt(),
	)
	require.NoError(t, err)

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeEncodingConfig().TxConfig
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, "", []uint64{0}, []uint64{0}, false, false, priv1)
	require.Error(t, err)
	require.EqualError(t, err, "commission can't be lower than 5%: unauthorized")
}

func TestCreateAndEditValidatorAnteHandler(t *testing.T) {
	genTokens := sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	bondTokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

	acc1 := &authtypes.BaseAccount{Address: addr1.String()}
	acc2 := &authtypes.BaseAccount{Address: addr2.String()}
	accs := authtypes.GenesisAccounts{acc1, acc2}
	balances := []banktypes.Balance{
		{
			Address: addr1.String(),
			Coins:   sdk.Coins{genCoin},
		},
		{
			Address: addr2.String(),
			Coins:   sdk.Coins{genCoin},
		},
	}

	app := simapp.SetupWithGenesisAccounts(accs, balances...)
	simapp.CheckBalance(t, app, addr1, sdk.Coins{genCoin})
	simapp.CheckBalance(t, app, addr2, sdk.Coins{genCoin})

	// create validator
	description := types.NewDescription("foo_moniker", "", "", "", "")
	createValidatorMsg, err := types.NewMsgCreateValidator(
		sdk.ValAddress(addr1), valKey.PubKey(), bondCoin, description, goodCommissionRates, sdk.OneInt(),
	)
	require.NoError(t, err)

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeEncodingConfig().TxConfig
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, "", []uint64{8}, []uint64{0}, true, true, priv1)
	require.NoError(t, err)
	simapp.CheckBalance(t, app, addr1, sdk.Coins{genCoin.Sub(bondCoin)})

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	validator := checkValidator(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, sdk.ValAddress(addr1).String(), validator.OperatorAddress)
	require.Equal(t, types.Bonded, validator.Status)
	require.True(sdk.IntEq(t, bondTokens, validator.BondedTokens()))

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// edit the validator
	description = types.NewDescription("bar_moniker", "", "", "", "")
	editValidatorMsg := types.NewMsgEditValidator(sdk.ValAddress(addr1), description, nil, nil)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{editValidatorMsg}, "", []uint64{8}, []uint64{1}, true, true, priv1)
	require.NoError(t, err)

	validator = checkValidator(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, description, validator.Description)

	// edit the validator - fail
	description = types.NewDescription("bar_moniker", "", "", "", "")
	zeroDec := sdk.ZeroDec()
	editValidatorMsg = types.NewMsgEditValidator(sdk.ValAddress(addr1), description, &zeroDec, nil)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{editValidatorMsg}, "", []uint64{8}, []uint64{1}, false, false, priv1)
	require.Error(t, err)
	require.EqualError(t, err, "commission can't be lower than 5%: unauthorized")
}

func TestMinCommissionAuthzAnteHandler(t *testing.T) {
	priv1 := secp256k1.GenPrivKey()
	pub1 := priv1.PubKey()
	addr1 := sdk.AccAddress(pub1.Address())

	priv2 := secp256k1.GenPrivKey()
	pub2 := priv2.PubKey()
	addr2 := sdk.AccAddress(pub2.Address())

	genTokens := sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	bondTokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

	acc1 := &authtypes.BaseAccount{Address: addr1.String()}
	acc2 := &authtypes.BaseAccount{Address: addr2.String()}
	accs := authtypes.GenesisAccounts{acc1, acc2}
	balances := []banktypes.Balance{
		{
			Address: addr1.String(),
			Coins:   sdk.Coins{genCoin},
		},
		{
			Address: addr2.String(),
			Coins:   sdk.Coins{genCoin},
		},
	}

	valKey := ed25519.GenPrivKey()
	// valAddr := sdk.AccAddress(valKey.PubKey().Address())

	commissionRates := types.NewCommissionRates(sdk.MustNewDecFromStr("0.04"), sdk.MustNewDecFromStr("0.10"), sdk.MustNewDecFromStr("0.01"))

	app := simapp.SetupWithGenesisAccounts(accs, balances...)

	key := valKey.PubKey()

	auth1 := authz.NewGenericAuthorization("/cosmos.staking.v1beta1.MsgCreateValidator")

	msg1, err := authz.NewMsgGrant(addr1, addr2, auth1, time.Now().Add(time.Hour*72))
	require.NotNil(t, msg1)
	require.NoError(t, err)

	auth2 := authz.NewGenericAuthorization("/cosmos.staking.v1beta1.MsgEditValidator")

	msg2, err := authz.NewMsgGrant(addr1, addr2, auth2, time.Now().Add(time.Hour*72))
	require.NoError(t, err)

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}

	encoding := simapp.MakeEncodingConfig()
	txGen := encoding.TxConfig

	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{msg1, msg2}, "",
		[]uint64{8}, []uint64{0}, true, true, priv1)
	require.NoError(t, err)

	// create 2 authorization

	// create validator
	description := types.NewDescription("foo_moniker", "", "", "", "")
	createValidatorMsg, err := types.NewMsgCreateValidator(
		sdk.ValAddress(addr1), key, bondCoin, description, commissionRates, sdk.OneInt(),
	)
	require.NotNil(t, createValidatorMsg)
	require.NoError(t, err)

	execMsg := authz.NewMsgExec(addr2, []sdk.Msg{createValidatorMsg})

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{&execMsg}, "",
		[]uint64{1}, []uint64{0}, false, false, priv2)
	require.EqualError(t, err, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "commission can't be lower than 5%").Error())

	// valid
	commissionRates = types.NewCommissionRates(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("0.20"), sdk.MustNewDecFromStr("0.1"))
	createValidatorMsg, err = types.NewMsgCreateValidator(
		sdk.ValAddress(addr1), key, bondCoin, description, commissionRates, sdk.OneInt(),
	)

	require.NoError(t, err)
	header = tmproto.Header{Height: app.LastBlockHeight() + 1}

	// wrapped tx
	execMsg = authz.NewMsgExec(addr2, []sdk.Msg{createValidatorMsg})
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{&execMsg}, "", []uint64{9}, []uint64{0}, true, true, priv2)
	require.NoError(t, err)
	validator := checkValidator(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, description, validator.Description)
	require.Equal(t, commissionRates.Rate, validator.Commission.Rate)

	// edit the validator to 1%
	description = types.NewDescription("low commission", "", "", "", "")
	com := sdk.MustNewDecFromStr("0.01")
	editValidatorMsg := types.NewMsgEditValidator(sdk.ValAddress(addr1), description, &com, nil)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}

	// wrapped tx
	execMsg = authz.NewMsgExec(addr2, []sdk.Msg{editValidatorMsg})

	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{&execMsg}, "", []uint64{1}, []uint64{1}, false, false, priv2)
	require.EqualError(t, err, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "commission can't be lower than 5%").Error())

	validator = checkValidator(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, commissionRates.Rate, validator.Commission.Rate)

	// edit the validator to 10%
	description = types.NewDescription("increase commission", "", "", "", "")
	com = sdk.MustNewDecFromStr("0.09")
	editValidatorMsg = types.NewMsgEditValidator(sdk.ValAddress(addr1), description, &com, nil)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1, Time: time.Now().Add(time.Hour * 25)}

	// wrapped tx
	execMsg = authz.NewMsgExec(addr2, []sdk.Msg{editValidatorMsg})
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{&execMsg}, "", []uint64{9}, []uint64{1}, false, true, priv2)
	require.NoError(t, err)

	validator = checkValidator(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, sdk.MustNewDecFromStr("0.09"), validator.Commission.Rate)
}

func checkValidator(t *testing.T, app *simapp.BitsongApp, addr sdk.ValAddress, expFound bool) types.Validator {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	validator, found := app.StakingKeeper.GetValidator(ctxCheck, addr)

	require.Equal(t, expFound, found)
	return validator
}
