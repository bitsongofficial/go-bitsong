package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBitsongdExport(t *testing.T) {
	db := db.NewMemDB()
	bapp := NewBitsongApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, "")
	setGenesis(bapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newBapp := NewBitsongApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, "")
	_, _, err := newBapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := db.NewMemDB()
	app := NewBitsongApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, "")

	for acc := range maccPerms {
		require.True(t, app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(gapp *GoBitsong) error {

	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	gapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	gapp.Commit()
	return nil
}
