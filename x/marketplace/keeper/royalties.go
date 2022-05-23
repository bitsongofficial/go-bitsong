package keeper

import (
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ProcessRoyalties(ctx sdk.Context, metadata nfttypes.Metadata, authority sdk.AccAddress, denom string, amount uint64) error {
	if metadata.Data.SellerFeeBasisPoints > 100 {
		return nfttypes.ErrInvalidSellerFeeBasisPoints
	}
	totalRoyalties := amount * uint64(metadata.Data.SellerFeeBasisPoints) / 100

	totalShare := uint32(0)
	for _, creator := range metadata.Data.Creators {
		totalShare += creator.Share
	}
	if totalShare == 0 {
		return nil
	}

	for _, creator := range metadata.Data.Creators {
		amount = totalRoyalties * uint64(creator.Share) / uint64(totalShare)
		if amount == 0 {
			continue
		}
		creatorAddr, err := sdk.AccAddressFromBech32(creator.Address)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoins(ctx, authority, creatorAddr, sdk.Coins{sdk.NewInt64Coin(denom, int64(amount))})
		if err != nil {
			return err
		}
	}
	return nil
}
