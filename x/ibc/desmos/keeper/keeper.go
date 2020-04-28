package keeper

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/ibc/desmos/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

const (
	// DefaultPacketTimeout is the default packet timeout relative to the current block height
	DefaultPacketTimeout = 1000 // NOTE: in blocks

	// DefaultPacketTimeoutTimestamp is the default packet timeout timestamp relative
	// to the current block timestamp. The timeout is disabled when set to 0.
	DefaultPacketTimeoutTimestamp = 0 // NOTE: in nanoseconds
)

// Represents the keeper that is used to perform IBC operations
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	scopedKeeper  capability.ScopedKeeper
}

func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey,
	ck types.ChannelKeeper, portK types.PortKeeper,
	sk capability.ScopedKeeper,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		channelKeeper: ck,
		portKeeper:    portK,
		scopedKeeper:  sk,
	}
}

// GetPort returns the portID for the IBC posts module.
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.PortKey)))
}

// SendPostCreation handles the creation of a post to a Desmos-based chain.
func (k Keeper) SendPostCreation(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	destHeight uint64,

	songID string,
	creationTime time.Time,
	postOwner string,
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
		songID, creationTime, postOwner,
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
	sender string,
) error {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packet := channel.NewPacket(
		types.NewSongCreationData(songID, creationTime, sender).GetBytes(),
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		destHeight+DefaultPacketTimeout,
		DefaultPacketTimeoutTimestamp,
	)

	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	// Set the portID into our store so we can retrieve it later
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PortKey), []byte(portID))

	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, porttypes.PortPath(portID))
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
