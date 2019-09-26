package distribution

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/BitSongOfficial/go-bitsong/x/track"

	tracktypes "github.com/BitSongOfficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type OverrideDistrKeeper struct {
	distr.Keeper
	stakingKeeper    types.StakingKeeper
	supplyKeeper     types.SupplyKeeper
	trackKeeper      track.Keeper
	feeCollectorName string
}

func NewOverrideDistrKeeper(keeper distr.Keeper, stakingKeeper types.StakingKeeper, supplyKeeper types.SupplyKeeper, trackKeeper track.Keeper, feeCollectorName string) OverrideDistrKeeper {
	return OverrideDistrKeeper{
		Keeper:           keeper,
		stakingKeeper:    stakingKeeper,
		supplyKeeper:     supplyKeeper,
		trackKeeper:      trackKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// AllocateTokens handles distribution of the collected fees
func (k OverrideDistrKeeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, previousVotes []abci.VoteInfo,
) {

	logger := k.Logger(ctx)

	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := feeCollector.GetCoins()
	feesCollected := sdk.NewDecCoins(feesCollectedInt)

	// transfer collected fees to the distribution module account
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected)
		k.SetFeePool(ctx, feePool)
		return
	}

	// calculate fraction votes
	previousFractionVotes := sdk.NewDec(sumPreviousPrecommitPower).Quo(sdk.NewDec(totalPreviousPower))

	// calculate previous proposer reward
	baseProposerReward := k.GetBaseProposerReward(ctx)
	bonusProposerReward := k.GetBonusProposerReward(ctx)
	proposerMultiplier := baseProposerReward.Add(bonusProposerReward.MulTruncate(previousFractionVotes))
	proposerReward := feesCollected.MulDecTruncate(proposerMultiplier)

	// pay previous proposer
	remaining := feesCollected
	proposerValidator := k.stakingKeeper.ValidatorByConsAddr(ctx, previousProposer)

	if proposerValidator != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetOperator().String()),
			),
		)

		k.AllocateTokensToValidator(ctx, proposerValidator, proposerReward)
		remaining = remaining.Sub(proposerReward)
	} else {
		// previous proposer can be unknown if say, the unbonding period is 1 block, so
		// e.g. a validator undelegates at block X, it's removed entirely by
		// block X+1's endblock, then X+2 we need to refer to the previous
		// proposer for X+1, but we've forgotten about them.
		logger.Error(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unbonded completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			previousProposer.String()))
	}

	// calculate fraction allocated to validators
	communityTax := k.GetCommunityTax(ctx)
	playTax := k.trackKeeper.GetPlayTax(ctx)
	voteMultiplier := sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax).Sub(playTax)

	// allocate tokens proportionally to voting power
	// TODO consider parallelizing later, ref https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
	for _, vote := range previousVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		// TODO consider microslashing for missing votes.
		// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
		powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPreviousPower))
		reward := feesCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}

	playPoolMultiplier := sdk.OneDec().Sub(communityTax).Sub(playTax)
	playPoolReward := remaining.MulDec(playPoolMultiplier)

	// truncate coins, return remainder to community pool
	coin, remainder := sdk.NewDecCoinFromDec("ubtsg", playPoolReward.AmountOf("ubtsg")).TruncateDecimal()

	// transfer collected play fees to the track module account
	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, tracktypes.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		panic(err)
	}

	playPool := k.trackKeeper.GetFeePlayPool(ctx)
	playPool.Rewards = playPool.Rewards.Add(sdk.NewDecCoins(sdk.NewCoins(coin)))
	k.trackKeeper.SetFeePlayPool(ctx, playPool)

	remaining = remaining.Sub(playPoolReward).Add(sdk.NewDecCoins(sdk.NewCoins(sdk.NewCoin(remainder.Denom, remainder.Amount.TruncateInt()))))

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining)
	k.SetFeePool(ctx, feePool)

	fmt.Println()
	fmt.Println()
	fmt.Printf("Play Rewards Pool: %s", playPool.Rewards.String())
	fmt.Println()
	fmt.Printf("Community Pool: %s", feePool.CommunityPool.String())
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
