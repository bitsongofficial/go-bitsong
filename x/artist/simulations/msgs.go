package simulations

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"math"
	"math/rand"

	"github.com/bitsongofficial/go-bitsong/x/artist"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// MetaSimulator defines a function type alias for generating random artists meta content.
type MetaSimulator func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) types.Meta

// SimulateCreatingArtist simulates creating a msg Create Artist on the artist module.
// It is implemented using future operations.
// TODO: improve simulations simulations
func SimulateCreatingArtist(k artist.Keeper, metaSim MetaSimulator) simulation.Operation {
	handler := artist.NewHandler(k)

	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account,
	) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// 1) create artist
		sender := simulation.RandomAcc(r, accs)
		meta := metaSim(r, app, ctx, accs)
		msg, err := simulationCreateMsgCreateArtist(r, meta, sender)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		ok := simulateHandleMsgCreateArtist(msg, handler, ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, meta.GetName())
		if !ok {
			return opMsg, nil, nil
		}

		artistID, err := k.GetArtistID(ctx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		artistID = uint64(math.Max(float64(artistID)-1, 0))

		y := 10
		fops := make([]simulation.FutureOperation, y)
		/*for i := 0; i < y; i++ {
			when := ctx.BlockHeader().Time
			fops[i] = simulation.FutureOperation{BlockTime: when, Op: operationSimulateMsgVote(k, accs[whoVotes[i]], proposalID)}
		}*/

		return opMsg, fops, nil
	}
}

func simulateHandleMsgCreateArtist(msg types.MsgCreateArtist, handler sdk.Handler, ctx sdk.Context) (ok bool) {
	ctx, write := ctx.CacheContext()
	ok = handler(ctx, msg).IsOK()
	if ok {
		write()
	}
	return ok
}

func simulationCreateMsgCreateArtist(r *rand.Rand, m types.Meta, s simulation.Account) (msg types.MsgCreateArtist, err error) {
	msg = types.NewMsgCreateArtist(m, s.Address)
	if msg.ValidateBasic() != nil {
		err = fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
	}
	return
}

// Pick a random artist ID
func randomArtistID(r *rand.Rand, k artist.Keeper, ctx sdk.Context) (artistID uint64, ok bool) {
	lastArtistID, _ := k.GetArtistID(ctx)
	lastArtistID = uint64(math.Max(float64(lastArtistID)-1, 0))

	if lastArtistID < 1 || lastArtistID == (2<<63-1) {
		return 0, false
	}
	artistID = uint64(r.Intn(1+int(lastArtistID)) - 1)
	return artistID, true
}
