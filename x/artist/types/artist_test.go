package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArtistStatus_Format(t *testing.T) {
	statusNil, _ := ArtistStatusFromString("")
	statusFailed, _ := ArtistStatusFromString("Failed")
	statusRejected, _ := ArtistStatusFromString("Rejected")
	statusVerified, _ := ArtistStatusFromString("Verified")
	tests := []struct {
		at                   ArtistStatus
		sprintFArgs          string
		expectedStringOutput string
	}{
		{statusNil, "%s", ""},
		{statusNil, "%v", "0"},
		{statusVerified, "%s", "Verified"},
		{statusVerified, "%v", "1"},
		{statusRejected, "%s", "Rejected"},
		{statusRejected, "%v", "2"},
		{statusFailed, "%s", "Failed"},
		{statusFailed, "%v", "3"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.sprintFArgs, tt.at)
		require.Equal(t, tt.expectedStringOutput, got)
	}
}
