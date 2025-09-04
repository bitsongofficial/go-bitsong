package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MintNFT(ctx sdk.Context, collectionDenom string, minter sdk.AccAddress, owner sdk.AccAddress, metadata types.Nft) error {
	if strings.TrimSpace(metadata.TokenId) == "" {
		return fmt.Errorf("token ID cannot be empty")
	}

	nftKey := collections.Join(collectionDenom, metadata.TokenId)
	has, err := k.NFTs.Has(ctx, nftKey)
	if err != nil {
		return fmt.Errorf("failed to check NFT: %w", err)
	}
	if has {
		return fmt.Errorf("NFT with token ID %s already exists in collection %s", metadata.TokenId, collectionDenom)
	}

	coll, err := k.Collections.Get(ctx, collectionDenom)
	if err != nil {
		return types.ErrCollectionNotFound
	}

	if coll.Minter == "" {
		return fmt.Errorf("minting disabled for this collection")
	}

	collectionMinter, err := sdk.AccAddressFromBech32(coll.Minter)
	if err != nil {
		return fmt.Errorf("invalid minter address: %w", err)
	}

	if !minter.Equals(collectionMinter) {
		return fmt.Errorf("only the collection minter can mint NFTs")
	}

	// TODO: Charge fee if necessary

	metadata.Collection = collectionDenom
	metadata.Owner = owner.String()

	if err := k.setNft(ctx, collectionDenom, metadata.TokenId, metadata); err != nil {
		return fmt.Errorf("failed to set NFT: %w", err)
	}

	return k.incrementSupply(ctx, collectionDenom)
}

func (k Keeper) createNftDenom(ctx sdk.Context, collectionDenom string) string {
	supply := k.GetSupply(ctx, collectionDenom)
	return fmt.Sprintf("%s-%d", collectionDenom, supply.Uint64()+1)
}

func (k Keeper) setNft(ctx sdk.Context, collectionDenom string, tokenId string, nft types.Nft) error {
	pk := collections.Join(collectionDenom, tokenId)
	return k.NFTs.Set(ctx, pk, nft)
}
