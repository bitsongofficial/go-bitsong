package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CreateCollection(
	ctx sdk.Context,
	creator,
	minter,
	authority,
	symbol,
	name,
	uri string,
) (denom string, err error) {
	creatorAddr, err := k.ac.StringToBytes(creator)
	if err != nil {
		return "", errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	denom, err = k.validateCollectionDenom(ctx, creatorAddr, symbol)
	if err != nil {
		return "", err
	}

	// TODO: charge fee

	if err := k.validateCollectionMetadata(name, uri); err != nil {
		return "", err
	}

	coll := types.Collection{
		Denom:   denom,
		Symbol:  symbol,
		Name:    name,
		Uri:     uri,
		Creator: creator,
	}

	if minter != "" {
		_, err = k.ac.StringToBytes(minter)
		if err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address: %s", err)
		}

		coll.Minter = minter
	}

	if authority != "" {
		_, err = k.ac.StringToBytes(authority)
		if err != nil {
			return "", errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %s", err)
		}

		coll.Authority = authority
	}

	if err := k.setCollection(ctx, coll); err != nil {
		return "", err
	}

	return denom, nil
}

func (k Keeper) SetCollectionName(ctx context.Context, authority sdk.AccAddress, denom, name string) error {
	coll, err := k.GetCollection(ctx, denom)
	if err != nil {
		return err
	}

	if coll.Authority != authority.String() {
		return fmt.Errorf("only the collection authority can change the name")
	}

	if err := k.validateCollectionMetadata(name, coll.Uri); err != nil {
		return err
	}

	coll.Name = name

	return k.setCollection(ctx, coll)
}

func (k Keeper) SetCollectionUri(ctx context.Context, authority sdk.AccAddress, denom, uri string) error {
	coll, err := k.GetCollection(ctx, denom)
	if err != nil {
		return err
	}

	if coll.Authority != authority.String() {
		return fmt.Errorf("only the collection authority can change the uri")
	}

	if err := k.validateCollectionMetadata(coll.Name, uri); err != nil {
		return err
	}

	coll.Uri = uri

	return k.setCollection(ctx, coll)
}

func (k Keeper) SetMinter(ctx context.Context, oldMinter sdk.AccAddress, newMinter *sdk.AccAddress, denom string) error {
	coll, err := k.GetCollection(ctx, denom)
	if err != nil {
		return err
	}

	if coll.Minter == "" {
		return fmt.Errorf("minting disabled for this collection")
	}

	if coll.Minter != oldMinter.String() {
		return fmt.Errorf("only the current minter can change the minter")
	}

	if newMinter != nil {
		coll.Minter = newMinter.String()
	} else {
		coll.Minter = ""
	}

	return k.setCollection(ctx, coll)
}

func (k Keeper) SetAuthority(ctx context.Context, oldAuthority sdk.AccAddress, newAuthority *sdk.AccAddress, denom string) error {
	coll, err := k.GetCollection(ctx, denom)
	if err != nil {
		return err
	}

	if coll.Authority != oldAuthority.String() {
		return fmt.Errorf("only the current authority can change the authority")
	}

	if newAuthority != nil {
		coll.Authority = newAuthority.String()
	} else {
		coll.Authority = ""
	}

	return k.setCollection(ctx, coll)
}

func (k Keeper) GetMinter(ctx context.Context, denom string) (sdk.AccAddress, error) {
	coll, err := k.Collections.Get(ctx, denom)
	if err != nil {
		return nil, types.ErrCollectionNotFound
	}

	if coll.Minter == "" {
		return nil, fmt.Errorf("minting disabled for this collection")
	}

	return sdk.AccAddressFromBech32(coll.Minter)
}

func (k Keeper) GetAuthority(ctx context.Context, denom string) (sdk.AccAddress, error) {
	coll, err := k.Collections.Get(ctx, denom)
	if err != nil {
		return nil, types.ErrCollectionNotFound
	}

	if coll.Authority == "" {
		return nil, fmt.Errorf("no authority set for this collection")
	}

	return sdk.AccAddressFromBech32(coll.Authority)
}

func (k Keeper) GetSupply(ctx context.Context, denom string) math.Int {
	supply, err := k.Supply.Get(ctx, denom)
	if err != nil {
		return math.ZeroInt()
	}

	return supply
}

func (k Keeper) HasSupply(ctx context.Context, denom string) bool {
	has, err := k.Supply.Has(ctx, denom)
	return has && err == nil
}

func (k Keeper) HasCollection(ctx context.Context, denom string) bool {
	has, err := k.Collections.Has(ctx, denom)
	return has && err == nil
}

func (k Keeper) GetCollection(ctx context.Context, denom string) (types.Collection, error) {
	coll, err := k.Collections.Get(ctx, denom)
	if err != nil {
		return types.Collection{}, types.ErrCollectionNotFound
	}

	return coll, nil
}

func (k Keeper) setSupply(ctx context.Context, denom string, supply math.Int) error {
	return k.Supply.Set(ctx, denom, supply)
}

func (k Keeper) incrementSupply(ctx context.Context, denom string) error {
	supply := k.GetSupply(ctx, denom)
	supply = supply.Add(math.NewInt(1))

	return k.setSupply(ctx, denom, supply)
}

func (k Keeper) createCollectionDenom(creator sdk.AccAddress, symbol string) (string, error) {
	// TODO: if necessary add a salt field

	if strings.TrimSpace(symbol) == "" {
		return "", fmt.Errorf("symbol cannot be empty")
	}

	if len(symbol) > types.MaxSymbolLength {
		return "", fmt.Errorf("symbol cannot be longer than %d characters", types.MaxSymbolLength)
	}

	bz := []byte(fmt.Sprintf("%s/%s", creator.String(), symbol))
	return fmt.Sprintf("nft%x", tmcrypto.AddressHash(bz)), nil
}

func (k Keeper) validateCollectionDenom(ctx context.Context, creator sdk.AccAddress, symbol string) (string, error) {
	denom, err := k.createCollectionDenom(creator, symbol)
	if err != nil {
		return "", err
	}

	if err := sdk.ValidateDenom(denom); err != nil {
		return "", err
	}

	if k.HasCollection(ctx, denom) {
		return "", types.ErrCollectionAlreadyExists
	}

	return denom, nil
}

func (k Keeper) setCollection(ctx context.Context, coll types.Collection) error {
	return k.Collections.Set(ctx, coll.Denom, coll)
}

func (k Keeper) validateCollectionMetadata(name, uri string) error {
	if len(name) > types.MaxNameLength {
		return fmt.Errorf("name cannot be longer than %d characters", types.MaxNameLength)
	}

	if len(uri) > types.MaxURILength {
		return fmt.Errorf("uri cannot be longer than %d characters", types.MaxURILength)
	}

	return nil
}
