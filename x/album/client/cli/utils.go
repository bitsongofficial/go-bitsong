package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"io/ioutil"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)

type (
	// CreateAlbumJSON defines a CreateAlbum msg
	CreateAlbumJSON struct {
		AlbumType   types.AlbumType `json:"album_type" yaml:"album_type"`
		Title       string          `json:"title" yaml:"title"`
		MetadataURI string          `json:"metadata_uri" yaml:"metadata_uri"`
	}
)

// ParseCreateAlbumJSON reads and parses a CreateAlbumJSON from a file.
func ParseCreateAlbumJSON(cdc *codec.Codec, albumFile string) (CreateAlbumJSON, error) {
	album := CreateAlbumJSON{}

	payload, err := ioutil.ReadFile(albumFile)
	if err != nil {
		return album, err
	}

	if err := cdc.UnmarshalJSON(payload, &album); err != nil {
		return album, err
	}

	return album, nil
}

func QueryDepositsByTxQuery(cliCtx context.CLIContext, params types.QueryAlbumParams) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgDeposit),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeDepositAlbum, types.AttributeKeyAlbumID, []byte(fmt.Sprintf("%d", params.AlbumID))),
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
					AlbumID:   params.AlbumID,
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
