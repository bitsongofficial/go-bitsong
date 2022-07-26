package fantoken

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// NewHandler handles all fantoken type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(&k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgIssue:
			res, err := msgServer.Issue(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDisableMint:
			res, err := msgServer.DisableMint(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgMint:
			res, err := msgServer.Mint(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBurn:
			res, err := msgServer.Burn(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSetAuthority:
			res, err := msgServer.SetAuthority(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSetMinter:
			res, err := msgServer.SetMinter(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSetUri:
			res, err := msgServer.SetUri(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UpdateFeesProposal:
			return handleUpdateFeesProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized fantoken proposal content type: %T", c)
		}
	}
}

func handleUpdateFeesProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateFeesProposal) error {
	ctx.Logger().Info("Updating fantoken fees from proposal")

	if err := types.ValidateFees(p.IssueFee, p.MintFee, p.BurnFee); err != nil {
		return err
	}

	params := k.GetParamSet(ctx)
	params.IssueFee = p.IssueFee
	params.MintFee = p.MintFee
	params.BurnFee = p.BurnFee
	k.SetParamSet(ctx, params)

	return nil
}
