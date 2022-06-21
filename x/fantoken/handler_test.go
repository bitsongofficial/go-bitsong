package fantoken_test

import (
	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/fantoken"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestProposalHandlerPassed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	params := app.FanTokenKeeper.GetParamSet(ctx)
	require.Equal(t, params, fantokentypes.DefaultParams())

	newIssueFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
	newMintFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2))
	newBurnFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3))
	newTransferFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))

	proposal := fantokentypes.NewUpdateFeesProposal(
		"Test",
		"description",
		newIssueFee,
		newMintFee,
		newBurnFee,
		newTransferFee,
	)

	h := fantoken.NewProposalHandler(app.FanTokenKeeper)
	require.NoError(t, h(ctx, proposal))

	params = app.FanTokenKeeper.GetParamSet(ctx)
	require.Equal(t, newIssueFee, params.IssueFee)
	require.Equal(t, newMintFee, params.MintFee)
	require.Equal(t, newBurnFee, params.BurnFee)
	require.Equal(t, newTransferFee, params.TransferFee)
}

func TestProposalHandlerFailed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	params := app.FanTokenKeeper.GetParamSet(ctx)
	require.Equal(t, params, fantokentypes.DefaultParams())

	newIssueFee := sdk.Coin{
		Denom:  sdk.DefaultBondDenom,
		Amount: sdk.NewInt(-1),
	}
	newMintFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())
	newBurnFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())
	newTransferFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())

	proposal := fantokentypes.NewUpdateFeesProposal(
		"Test",
		"description",
		newIssueFee,
		newMintFee,
		newBurnFee,
		newTransferFee,
	)

	h := fantoken.NewProposalHandler(app.FanTokenKeeper)
	require.Error(t, h(ctx, proposal))
}
