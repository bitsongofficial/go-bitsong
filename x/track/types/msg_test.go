package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
	"time"
)

var (
	mockTitle          = "The Show Must Go On"
	mockEmptyTitle     = ""
	mockEmptyAttribute = map[string]string{}
	mockContent        = Content{
		Path:        "/ipfs/Qm....",
		ContentType: "",
		Duration:    0,
		Attributes:  nil,
	}
	mockEmptyContent      = Content{}
	mockTrackMediaNoVideo = TrackMedia{
		Audio: mockContent,
		Video: mockEmptyContent,
		Image: mockContent,
	}
	mockEmptyTrackMedia     = TrackMedia{}
	mockRightHolder1        = NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder1"))), 100)
	mockRightHolder2        = NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder2"))), 25)
	mockRightHolder3        = NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder3"))), 25)
	mockRightHolder4        = NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder4"))), 50)
	mockRightsHoldersSingle = RightsHolders{
		mockRightHolder1,
	}
	mockRightsHoldersMultiple = RightsHolders{
		mockRightHolder2,
		mockRightHolder3,
		mockRightHolder4,
	}
	mockRewards = TrackRewards{
		Users:     10,
		Playlists: 10,
	}
	mockOwner = sdk.AccAddress(crypto.AddressHash([]byte("owner")))
	mockTrack = Track{
		Title:         mockTitle,
		Attributes:    mockEmptyAttribute,
		Media:         mockTrackMediaNoVideo,
		Rewards:       mockRewards,
		RightsHolders: mockRightsHoldersSingle,
		SubmitTime:    time.Time{},
		Owner:         mockOwner,
	}
	mockTrackOwnerNil = Track{
		Title:         mockTitle,
		Attributes:    mockEmptyAttribute,
		Media:         mockTrackMediaNoVideo,
		Rewards:       mockRewards,
		RightsHolders: mockRightsHoldersSingle,
		SubmitTime:    time.Time{},
		Owner:         nil,
	}
)

var mockMsgCreate = NewMsgCreate(
	mockTitle,
	mockEmptyAttribute,
	mockTrackMediaNoVideo,
	mockRewards,
	mockRightsHoldersSingle,
	mockOwner,
)

func TestMsgCreate_Route(t *testing.T) {
	expected := "track"
	actual := mockMsgCreate.Route()
	require.Equal(t, expected, actual)
}

func TestMsgCreate_Type(t *testing.T) {
	expected := "create"
	actual := mockMsgCreate.Type()
	require.Equal(t, expected, actual)
}

func TestMsgCreate_ValidateBasic(t *testing.T) {
	_ = sdk.AccAddress(crypto.AddressHash([]byte("test")))

	// TODO: continue with more test
	tests := []struct {
		name  string
		msg   MsgCreate
		error error
	}{
		{
			name:  "Empty owner return error",
			msg:   NewMsgCreate(mockTitle, mockEmptyAttribute, mockTrackMediaNoVideo, mockRewards, mockRightsHoldersSingle, nil),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid owner: "),
		},
		{
			name:  "Empty title returns error if title is empty",
			msg:   NewMsgCreate(mockEmptyTitle, mockEmptyAttribute, mockTrackMediaNoVideo, mockRewards, mockRightsHoldersSingle, mockOwner),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "title cannot be blank"),
		},
		{
			name:  "Non-Empty title returns error if owner and media are empty",
			msg:   NewMsgCreate(mockTitle, mockEmptyAttribute, mockEmptyTrackMedia, mockRewards, mockRightsHoldersSingle, nil),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track audio is not a valid format"),
		},
		{
			name:  "Empty media returns error if media are empty",
			msg:   NewMsgCreate(mockTitle, mockEmptyAttribute, mockEmptyTrackMedia, mockRewards, mockRightsHoldersSingle, mockOwner),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track audio is not a valid format"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			returnedError := test.msg.ValidateBasic()
			if test.error == nil {
				require.Nil(t, returnedError)
			} else {
				require.NotNil(t, returnedError)
				require.Equal(t, test.error.Error(), returnedError.Error())
			}
		})
	}

	err := mockMsgCreate.ValidateBasic()
	require.Nil(t, err)
}

func TestMsgCreate_GetSignBytes(t *testing.T) {
	tests := []struct {
		name        string
		msg         MsgCreate
		expSignJSON string
	}{
		{
			name:        "Message with no attributes",
			msg:         NewMsgCreate(mockTitle, nil, mockTrackMediaNoVideo, mockRewards, mockRightsHoldersSingle, mockOwner),
			expSignJSON: `{"type":"go-bitsong/MsgCreateTrack","value":{"media":{"audio":{"attributes":null,"content_type":"","duration":0,"path":"/ipfs/Qm...."},"image":{"attributes":null,"content_type":"","duration":0,"path":"/ipfs/Qm...."}},"owner":"cosmos1fsgzj6t7udv8zhf6zj32mkqhcjcpv52ygswxa5","rewards":{"playlists":"10","users":"10"},"rights_holders":[{"address":"cosmos17zfanegzaj8shhzsrfncz6cz5ykvzr06yyww88","quota":"100"}],"title":"The Show Must Go On"}}`,
		},
		{
			name:        "Message with no attributes and no media",
			msg:         NewMsgCreate(mockTitle, nil, mockEmptyTrackMedia, mockRewards, mockRightsHoldersSingle, mockOwner),
			expSignJSON: `{"type":"go-bitsong/MsgCreateTrack","value":{"owner":"cosmos1fsgzj6t7udv8zhf6zj32mkqhcjcpv52ygswxa5","rewards":{"playlists":"10","users":"10"},"rights_holders":[{"address":"cosmos17zfanegzaj8shhzsrfncz6cz5ykvzr06yyww88","quota":"100"}],"title":"The Show Must Go On"}}`,
		},
		{
			name:        "Message with empty attributes and no media",
			msg:         NewMsgCreate(mockTitle, map[string]string{}, mockEmptyTrackMedia, mockRewards, mockRightsHoldersSingle, mockOwner),
			expSignJSON: `{"type":"go-bitsong/MsgCreateTrack","value":{"attributes":{},"owner":"cosmos1fsgzj6t7udv8zhf6zj32mkqhcjcpv52ygswxa5","rewards":{"playlists":"10","users":"10"},"rights_holders":[{"address":"cosmos17zfanegzaj8shhzsrfncz6cz5ykvzr06yyww88","quota":"100"}],"title":"The Show Must Go On"}}`,
		},
		{
			name:        "Message with empty attributes and media",
			msg:         NewMsgCreate(mockTitle, map[string]string{}, mockTrackMediaNoVideo, mockRewards, mockRightsHoldersSingle, mockOwner),
			expSignJSON: `{"type":"go-bitsong/MsgCreateTrack","value":{"attributes":{},"media":{"audio":{"attributes":null,"content_type":"","duration":0,"path":"/ipfs/Qm...."},"image":{"attributes":null,"content_type":"","duration":0,"path":"/ipfs/Qm...."}},"owner":"cosmos1fsgzj6t7udv8zhf6zj32mkqhcjcpv52ygswxa5","rewards":{"playlists":"10","users":"10"},"rights_holders":[{"address":"cosmos17zfanegzaj8shhzsrfncz6cz5ykvzr06yyww88","quota":"100"}],"title":"The Show Must Go On"}}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expSignJSON, string(test.msg.GetSignBytes()))
		})
	}
}

func TestMsgCreate_GetSigners(t *testing.T) {
	expected := mockMsgCreate.Owner
	actual := mockMsgCreate.GetSigners()
	require.Equal(t, expected, actual[0])
	require.Equal(t, 1, len(actual))
}
