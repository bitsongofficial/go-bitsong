package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k Keeper) CreateCollection(ctx sdk.Context, creator sdk.AccAddress, coll types.Collection) (denom string, err error) {
	denom, err = k.validateCollectionDenom(ctx, creator, coll.Symbol)
	if err != nil {
		return "", err
	}

	// TODO: charge fee

	metadata := banktypes.Metadata{
		DenomUnits: []*banktypes.DenomUnit{{
			Denom:    denom,
			Exponent: 0,
		}},
		Base:        denom,
		Name:        coll.Name,
		Description: coll.Description,
		Symbol:      coll.Symbol,
		Display:     coll.Symbol,
		URI:         coll.Uri,
	}

	k.bk.SetDenomMetaData(ctx, metadata)

	if err := k.setCollection(ctx, denom, coll); err != nil {
		return "", err
	}

	return denom, nil
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) math.Int {
	supply, err := k.Supply.Get(ctx, denom)
	if err != nil {
		return math.ZeroInt()
	}

	return supply
}

func (k Keeper) HasSupply(ctx sdk.Context, denom string) bool {
	has, err := k.Supply.Has(ctx, denom)
	return has && err == nil
}

func (k Keeper) setSupply(ctx sdk.Context, denom string, supply math.Int) error {
	return k.Supply.Set(ctx, denom, supply)
}

func (k Keeper) incrementSupply(ctx sdk.Context, denom string) error {
	supply := k.GetSupply(ctx, denom)
	supply = supply.Add(math.NewInt(1))

	return k.setSupply(ctx, denom, supply)
}

func (k Keeper) createCollectionDenom(creator sdk.AccAddress, symbol string) string {
	// TODO: if necessary add a salt field

	bz := []byte(fmt.Sprintf("%s/%s", creator.String(), symbol))
	return fmt.Sprintf("nft%x", tmcrypto.AddressHash(bz))
}

func (k Keeper) validateCollectionDenom(ctx sdk.Context, creator sdk.AccAddress, symbol string) (string, error) {
	denom := k.createCollectionDenom(creator, symbol)

	if err := sdk.ValidateDenom(denom); err != nil {
		return "", err
	}

	if k.bk.HasSupply(ctx, denom) {
		return "", fmt.Errorf("denom %s already exists", denom)
	}

	_, exists := k.bk.GetDenomMetaData(ctx, denom)
	if exists {
		return "", types.ErrCollectionAlreadyExists
	}

	return denom, nil
}

func (k Keeper) setCollection(ctx sdk.Context, denom string, coll types.Collection) error {
	return k.Collections.Set(ctx, denom, coll)
}

func (k Keeper) getCollection(ctx sdk.Context, denom string) (types.Collection, error) {
	coll, err := k.Collections.Get(ctx, denom)
	if err != nil {
		return types.Collection{}, types.ErrCollectionNotFound
	}

	return coll, nil
}
