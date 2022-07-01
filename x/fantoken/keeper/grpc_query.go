package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) FanToken(c context.Context, req *types.QueryFanTokenRequest) (*types.QueryFanTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if len(req.Denom) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty denom")
	}

	fantoken, err := k.GetFanToken(ctx, req.Denom)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "fan token %s not found", req.Denom)
	}

	return &types.QueryFanTokenResponse{Fantoken: fantoken}, nil
}

func (k Keeper) FanTokens(c context.Context, req *types.QueryFanTokensRequest) (*types.QueryFanTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var owner sdk.AccAddress
	var err error

	if len(req.Authority) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Authority)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid authority address (%s)", err))
		}
	}

	var fantokens []*types.FanToken
	var pageRes *query.PageResponse

	store := ctx.KVStore(k.storeKey)

	if owner == nil {
		fantokenStore := prefix.NewStore(store, types.PrefixFanTokenForDenom)

		pageRes, err = query.Paginate(fantokenStore, req.Pagination, func(key []byte, value []byte) error {
			var fantoken types.FanToken
			k.cdc.MustUnmarshal(value, &fantoken)
			fantokens = append(fantokens, &fantoken)
			return nil
		})

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	} else {
		fantokenStore := prefix.NewStore(store, types.KeyFanTokens(owner, ""))

		pageRes, err = query.Paginate(fantokenStore, req.Pagination, func(key []byte, value []byte) error {
			var denom gogotypes.StringValue
			k.cdc.MustUnmarshal(value, &denom)
			fantoken, err := k.GetFanToken(ctx, denom.Value)
			if err == nil {
				fantokens = append(fantokens, fantoken)
			}
			return nil
		})

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	}

	var result []*types.FanToken
	for _, fantoken := range fantokens {
		result = append(result, fantoken)
	}

	return &types.QueryFanTokensResponse{Fantokens: result, Pagination: pageRes}, nil
}

// Params return the all the parameter in fantoken module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetParamSet(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
