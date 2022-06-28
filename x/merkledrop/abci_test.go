package merkledrop_test

import (
	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/keeper"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestExpiredMerkledrop(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	msgSvr := keeper.NewMsgServerImpl(app.MerkledropKeeper)

	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	initCoins := sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000)),
	}
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, owner, initCoins)
	require.NoError(t, err)

	merkledrops := app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)

	startHeight := int64(0)
	endHeight := int64(100)

	airdropCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))

	msg := types.NewMsgCreate(
		owner,
		"",
		startHeight,
		endHeight,
		airdropCoin,
	)

	balance := app.BankKeeper.GetBalance(ctx, owner, sdk.DefaultBondDenom)
	require.Equal(t, initCoins.AmountOf(sdk.DefaultBondDenom), balance.Amount)

	res, err := msgSvr.Create(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	balance = app.BankKeeper.GetBalance(ctx, owner, sdk.DefaultBondDenom)
	creationFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))
	expectedAmt := initCoins.AmountOf(sdk.DefaultBondDenom).Sub(creationFee.Amount).Sub(airdropCoin.Amount)
	require.Equal(t, expectedAmt, balance.Amount)

	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 1)

	newHeader := ctx.BlockHeader()
	newHeader.Height = endHeight
	ctx = ctx.WithBlockHeader(newHeader)

	merkledrop.EndBlocker(ctx, app.MerkledropKeeper)

	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)

	balance = app.BankKeeper.GetBalance(ctx, owner, sdk.DefaultBondDenom)
	expectedAmt = initCoins.AmountOf(sdk.DefaultBondDenom).Sub(creationFee.Amount)
	require.Equal(t, expectedAmt, balance.Amount)
}
