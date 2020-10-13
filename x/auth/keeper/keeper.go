package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/auth/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type Keeper struct {
	authKeeper auth.AccountKeeper
}

func NewKeeper(ak auth.AccountKeeper) Keeper {
	return Keeper{authKeeper: ak}
}

func (k Keeper) RegisterHandle(ctx sdk.Context, addr sdk.AccAddress, handle string) (*types.BitSongAccount, error) {
	acc := k.authKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	base := auth.BaseAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		PubKey:        acc.GetPubKey(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	bacc := types.NewBitSongAccount(base, handle)
	k.authKeeper.SetAccount(ctx, bacc)

	return bacc, nil
}
