package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueToken         = "op_weight_msg_issue_token"
	OpWeightMsgEditToken          = "op_weight_msg_edit_token"
	OpWeightMsgMintToken          = "op_weight_msg_mint_token"
	OpWeightMsgTransferTokenOwner = "op_weight_msg_transfer_token_owner"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak tokentypes.AccountKeeper,
	bk tokentypes.BankKeeper,
) simulation.WeightedOperations {

	var weightIssue, weightEdit, weightMint, weightTransfer int
	appParams.GetOrGenerate(
		cdc, OpWeightMsgIssueToken, &weightIssue, nil,
		func(_ *rand.Rand) {
			weightIssue = 100
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditToken, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintToken, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferTokenOwner, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightIssue,
			SimulateIssueFanToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateEditFanToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMintFanToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateTransferFanTokenOwner(k, ak, bk),
		),
	}
}

// SimulateIssueToken tests and runs a single msg issue a new token
func SimulateIssueFanToken(k keeper.Keeper, ak tokentypes.AccountKeeper, bk tokentypes.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, maxFees := genFanToken(ctx, r, k, ak, bk, accs)
		msg := tokentypes.NewMsgIssueFanToken(token.GetSymbol(), token.Name, token.MaxSupply, token.MetaData.Description, token.GetOwner().String(), token.GetUri(), sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)))

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), fmt.Sprintf("account %s not found", token.Owner)), nil, fmt.Errorf("account %s not found", token.Owner)
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		fees, err := simtypes.RandomFees(r, ctx, maxFees)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate issue token", nil), nil, nil
	}
}

// SimulateEditToken tests and runs a single msg edit a existed token
func SimulateEditFanToken(k keeper.Keeper, ak tokentypes.AccountKeeper, bk tokentypes.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, _, skip := selectOneFanToken(ctx, k, ak, bk, false)
		if skip {
			return simtypes.NoOpMsg(tokentypes.ModuleName, tokentypes.TypeMsgEditFanToken, "skip edit token"), nil, nil
		}
		msg := tokentypes.NewMsgEditFanToken(token.GetSymbol(), true, token.GetOwner().String())

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), fmt.Sprintf("account %s not found", token.GetOwner())), nil, fmt.Errorf("account %s not found", token.GetOwner())
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate edit token", nil), nil, nil
	}
}

// SimulateMintToken tests and runs a single msg mint a existed token
func SimulateMintFanToken(k keeper.Keeper, ak tokentypes.AccountKeeper, bk tokentypes.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, maxFee, skip := selectOneFanToken(ctx, k, ak, bk, true)
		if skip {
			return simtypes.NoOpMsg(tokentypes.ModuleName, tokentypes.TypeMsgMintFanToken, "skip mint token"), nil, nil
		}
		simToAccount, _ := simtypes.RandomAcc(r, accs)
		msg := tokentypes.NewMsgMintFanToken(simToAccount.Address.String(), token.GetDenom(), token.GetOwner().String(), sdk.NewInt(100))

		ownerAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), fmt.Sprintf("account %s not found", token.GetOwner())), nil, fmt.Errorf("account %s not found", token.GetOwner())
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		fees, err := simtypes.RandomFees(r, ctx, maxFee)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			ownerAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate mint token", nil), nil, nil
	}
}

// SimulateTransferTokenOwner tests and runs a single msg transfer to others
func SimulateTransferFanTokenOwner(k keeper.Keeper, ak tokentypes.AccountKeeper, bk tokentypes.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, _, skip := selectOneFanToken(ctx, k, ak, bk, false)
		if skip {
			return simtypes.NoOpMsg(tokentypes.ModuleName, tokentypes.TypeMsgTransferFanTokenOwner, "skip TransferTokenOwner"), nil, nil
		}
		var simToAccount, _ = simtypes.RandomAcc(r, accs)
		for simToAccount.Address.Equals(token.GetOwner()) {
			simToAccount, _ = simtypes.RandomAcc(r, accs)
		}

		msg := tokentypes.NewMsgTransferFanTokenOwner(token.GetSymbol(), token.GetOwner().String(), simToAccount.Address.String())

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), fmt.Sprintf("account %s not found", token.GetOwner())), nil, fmt.Errorf("account %s not found", token.GetOwner())
		}

		srcOwner, _ := sdk.AccAddressFromBech32(msg.SrcOwner)
		account := ak.GetAccount(ctx, srcOwner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(tokentypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate transfer token", nil), nil, nil
	}
}

func selectOneFanToken(
	ctx sdk.Context,
	k keeper.Keeper,
	ak tokentypes.AccountKeeper,
	bk tokentypes.BankKeeper,
	mint bool,
) (token tokentypes.FanTokenI, maxFees sdk.Coins, skip bool) {
	tokens := k.GetFanTokens(ctx, nil)
	if len(tokens) == 0 {
		return token, maxFees, true
	}

	for _, t := range tokens {
		if !mint {
			return t, nil, false
		}

		account := ak.GetAccount(ctx, t.GetOwner())
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		spendableStake := spendable.AmountOf(sdk.DefaultBondDenom)
		if spendableStake.IsZero() {
			continue
		}
		maxFees = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, spendableStake))
		token = t
		return
	}
	return token, maxFees, true
}

func randStringBetween(r *rand.Rand, min, max int) string {
	strLen := simtypes.RandIntBetween(r, min, max)
	randStr := simtypes.RandStringOfLength(r, strLen)
	return strings.ToLower(randStr)
}

func genFanToken(ctx sdk.Context,
	r *rand.Rand,
	k keeper.Keeper,
	ak tokentypes.AccountKeeper,
	bk tokentypes.BankKeeper,
	accs []simtypes.Account,
) (tokentypes.FanToken, sdk.Coins) {

	var token tokentypes.FanToken
	token = randFanToken(r, accs)

	for k.HasFanToken(ctx, token.GetSymbol()) {
		token = randFanToken(r, accs)
	}

	issueFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000))

	account, maxFees := filterAccount(ctx, r, ak, bk, accs, issueFee)
	token.Owner = account.String()

	return token, maxFees
}

func filterAccount(
	ctx sdk.Context,
	r *rand.Rand,
	ak tokentypes.AccountKeeper,
	bk tokentypes.BankKeeper,
	accs []simtypes.Account, fee sdk.Coin,
) (owner sdk.AccAddress, maxFees sdk.Coins) {
loop:
	simAccount, _ := simtypes.RandomAcc(r, accs)
	account := ak.GetAccount(ctx, simAccount.Address)
	spendable := bk.SpendableCoins(ctx, account.GetAddress())
	spendableStake := spendable.AmountOf(sdk.DefaultBondDenom)
	if spendableStake.IsZero() || spendableStake.LT(fee.Amount) {
		goto loop
	}
	owner = account.GetAddress()
	maxFees = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, spendableStake).Sub(fee))
	return
}

func randFanToken(r *rand.Rand, accs []simtypes.Account) tokentypes.FanToken {
	symbol := randStringBetween(r, tokentypes.MinimumSymbolLen, tokentypes.MaximumSymbolLen)
	denom := fmt.Sprintf("%s%s", "u", symbol)
	name := randStringBetween(r, 1, tokentypes.MaximumNameLen)
	maxSupply := sdk.NewInt(10000000000)
	uri := randStringBetween(r, 0, tokentypes.MaximumUriLen)
	simAccount, _ := simtypes.RandomAcc(r, accs)

	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        denom,
		Display:     symbol,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: denom, Exponent: 0},
			{Denom: symbol, Exponent: tokentypes.FanTokenDecimal},
		},
	}

	return tokentypes.FanToken{
		Name:      name,
		MaxSupply: maxSupply,
		Mintable:  true,
		Owner:     simAccount.Address.String(),
		URI:       uri,
		MetaData:  denomMetaData,
	}
}
