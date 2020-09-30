package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/account/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey sdk.StoreKey
	codec    *codec.Codec

	accountKeeper auth.AccountKeeper
	//bankKeeper    BankKeeper
	//supplyKeeper  supply.Keeper
}

func NewKeeper(storeKey sdk.StoreKey, codec *codec.Codec, accountKeeper auth.AccountKeeper) Keeper {
	keeper := Keeper{
		storeKey:      storeKey,
		codec:         codec,
		accountKeeper: accountKeeper,
	}

	// ensure track module account is set
	/*if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}*/

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) setAccount(ctx sdk.Context, acc types.Account) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(&acc)
	store.Set(types.GetAccountKey(acc.Address), bz)
}

func (k Keeper) CreateAppAccount(ctx sdk.Context, address sdk.AccAddress, handle string) (account types.Account, err error) {

	// first create a base account
	baseAccount := auth.NewBaseAccountWithAddress(address)
	//err := baseAccount.SetPubKey(pubKey)
	/*if err != nil {
		return appAccnt, ErrAppAccountCreateFailed(address)
	}*/
	k.accountKeeper.SetAccount(ctx, &baseAccount)

	//  then create an app account
	account = types.NewAccount(handle, "", ctx.BlockHeader().Time)
	k.setAccount(ctx, account)

	k.Logger(ctx).Info(fmt.Sprintf("Created %s", account.String()))

	return account, nil
}
