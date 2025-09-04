package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateCollection(
	ctx context.Context,
	creator,
	minter sdk.AccAddress,
	symbol,
	name,
	description,
	uri string,
) (denom string, err error) {
	denom, err = k.validateCollectionDenom(ctx, creator, symbol)
	if err != nil {
		return "", err
	}

	// TODO: charge fee

	if err := k.validateCollectionMetadata(name, description, uri); err != nil {
		return "", err
	}

	coll := types.Collection{
		Denom:       denom,
		Symbol:      symbol,
		Name:        name,
		Description: description,
		Uri:         uri,
		Creator:     creator.String(),
		Minter:      minter.String(),
	}

	if err := k.setCollection(ctx, coll); err != nil {
		return "", err
	}

	return denom, nil
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
		return "", fmt.Errorf("symbol cannot be blank")
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

func (k Keeper) getCollection(ctx context.Context, denom string) (types.Collection, error) {
	coll, err := k.Collections.Get(ctx, denom)
	if err != nil {
		return types.Collection{}, types.ErrCollectionNotFound
	}

	return coll, nil
}

func (k Keeper) validateCollectionMetadata(name, description, uri string) error {
	if len(name) > types.MaxNameLength {
		return fmt.Errorf("name cannot be longer than %d characters", types.MaxNameLength)
	}

	if len(description) > types.MaxDescriptionLength {
		return fmt.Errorf("description cannot be longer than %d characters", types.MaxDescriptionLength)
	}

	if len(uri) > types.MaxURILength {
		return fmt.Errorf("uri cannot be longer than %d characters", types.MaxURILength)
	}

	return nil
}
