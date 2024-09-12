package keeper

import (
	"fmt"

	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Codec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper

	paramSpace types.ParamSubspace
}

func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	dk types.DistrKeeper,
	paramSpace paramstypes.Subspace,
) Keeper {
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the " + types.ModuleName + " module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: ak,
		bankKeeper:    bk,
		distrKeeper:   dk,
		paramSpace:    paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("go-bitsong/%s", types.ModuleName))
}

func (k Keeper) Withdraw(ctx sdk.Context, merkledropID uint64) error {
	// get merkledrop
	merkledrop, err := k.getMerkleDropById(ctx, merkledropID)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrMerkledropNotExist, "merkledrop: %d does not exist", merkledropID)
	}

	// check if total amount < claimed amount  (who knows?)
	if merkledrop.Amount.LT(merkledrop.Claimed) {
		panic(fmt.Errorf("merkledrop-id: %d, total_amount (%s) < claimed_amount (%s)", merkledrop.Id, merkledrop.Amount, merkledrop.Claimed))
	}

	// get balance
	balance := merkledrop.Amount.Sub(merkledrop.Claimed)

	// send coins
	coin := sdk.NewCoin(merkledrop.Denom, balance)
	owner, err := sdk.AccAddressFromBech32(merkledrop.Owner)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.Coins{coin})
	if err != nil {
		return sdkerrors.Wrapf(types.ErrTransferCoins, "%s", coin)
	}

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
		MerkledropId: merkledrop.Id,
		Coin:         coin,
	})

	return nil
}
