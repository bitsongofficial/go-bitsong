package keeper

import (
	"fmt"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k Keeper) MintNFT(ctx sdk.Context, collectionDenom string, minter sdk.AccAddress, owner sdk.AccAddress, metadata types.Nft) (string, error) {
	coll, err := k.Collections.Get(ctx, collectionDenom)
	if err != nil {
		return "", types.ErrCollectionNotFound
	}

	collectionMinter, err := sdk.AccAddressFromBech32(coll.Minter)
	if err != nil {
		return "", fmt.Errorf("invalid minter address: %w", err)
	}

	if !minter.Equals(collectionMinter) {
		return "", fmt.Errorf("only the collection minter can mint NFTs")
	}

	nftDenom := k.createNftDenom(ctx, collectionDenom)

	// TODO: Charge fee if necessary

	nftMetadata := banktypes.Metadata{
		DenomUnits: []*banktypes.DenomUnit{{
			Denom:    nftDenom,
			Exponent: 0,
		}},
		Base:        nftDenom,
		Name:        metadata.Name,
		Description: metadata.Description,
		URI:         metadata.Uri,
		Symbol:      nftDenom,
		Display:     nftDenom,
	}

	k.bk.SetDenomMetaData(ctx, nftMetadata)

	amount := sdk.NewInt64Coin(nftDenom, 1)

	if err := k.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return "", fmt.Errorf("failed to mint NFT: %w", err)
	}

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(amount)); err != nil {
		return "", fmt.Errorf("failed to send NFT to owner: %w", err)
	}

	if err := k.incrementSupply(ctx, collectionDenom); err != nil {
		return "", fmt.Errorf("failed to increment supply: %w", err)
	}

	return nftDenom, nil
}

func (k Keeper) createNftDenom(ctx sdk.Context, collectionDenom string) string {
	supply := k.GetSupply(ctx, collectionDenom)
	return fmt.Sprintf("%s-%d", collectionDenom, supply.Uint64()+1)
}
