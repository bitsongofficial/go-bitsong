package keeper

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/ibc/desmos/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	ibcxfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	"github.com/desmos-labs/desmos/x/posts"
)

const (
	// DefaultPacketTimeout is the default packet timeout relative to the current block height
	DefaultPacketTimeout = 1000 // NOTE: in blocks
)

// Represents the keeper that is used to perform IBC operations
type Keeper struct {
	cdc           *codec.Codec
	channelKeeper ibcxfer.ChannelKeeper
}

func NewKeeper(cdc *codec.Codec, ck ibcxfer.ChannelKeeper) Keeper {
	return Keeper{
		cdc:           cdc,
		channelKeeper: ck,
	}
}

// SendPostCreation handles the creation of a post to a Desmos-based chain.
func (k Keeper) SendPostCreation(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	destHeight uint64,

	songID string,
	creationTime time.Time,
	sender sdk.AccAddress,
) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channel.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.Counterparty.PortID
	destinationChannel := sourceChannelEnd.Counterparty.ChannelID

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channel.ErrSequenceSendNotFound
	}

	return k.createOutgoingPacket(
		ctx, sequence, sourcePort, sourceChannel, destinationPort, destinationChannel, destHeight,
		songID, creationTime, sender,
	)
}

func (k Keeper) createOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	sourcePort, sourceChannel,
	destinationPort, destinationChannel string,
	destHeight uint64,

	songID string,
	creationTime time.Time,
	sender sdk.AccAddress,
) error {
	packetData := posts.NewCreatePostPacketData(
		types.NewSongCreationData(songID, creationTime, sender),
	)

	packet := channel.NewPacket(
		packetData.GetBytes(),
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		destHeight+DefaultPacketTimeout,
	)

	return k.channelKeeper.SendPacket(ctx, packet)
}
