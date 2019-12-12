package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)

func QueryArtistByID(artistID uint64, cliCtx context.CLIContext, queryRoute string) ([]byte, error) {
	params := types.NewQueryArtistParams(artistID)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/artist", queryRoute), bz)
	if err != nil {
		return nil, err
	}

	return res, err
}

// QueryDepositsByTxQuery will query for deposits via a direct txs tags query. It
// will fetch and build deposits directly from the returned txs and return a
// JSON marshalled result or any error that occurred.
//
// NOTE: SearchTxs is used to facilitate the txs query which does not currently
// support configurable pagination.
func QueryDepositsByTxQuery(cliCtx context.CLIContext, params types.QueryArtistParams) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgDeposit),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeDepositArtist, types.AttributeKeyArtistID, []byte(fmt.Sprintf("%d", params.ArtistID))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := utils.QueryTxsByEvents(cliCtx, events, defaultPage, defaultLimit)
	if err != nil {
		return nil, err
	}

	var deposits []types.Deposit

	for _, info := range searchResult.Txs {
		for _, msg := range info.Tx.GetMsgs() {
			if msg.Type() == types.TypeMsgDeposit {
				depMsg := msg.(types.MsgDeposit)

				deposits = append(deposits, types.Deposit{
					Depositor: depMsg.Depositor,
					ArtistID:  params.ArtistID,
					Amount:    depMsg.Amount,
				})
			}
		}
	}

	if cliCtx.Indent {
		return cliCtx.Codec.MarshalJSONIndent(deposits, "", "  ")
	}

	return cliCtx.Codec.MarshalJSON(deposits)
}
