package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	//govutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	artistTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Artist transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	artistTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateArtist(cdc),
	)...)

	return artistTxCmd
}

// GetCmdCreateArtist implements the create artist command handler.
func GetCmdCreateArtist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-artist",
		Short: "create new artist initialized with status nil",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Artist initialized with status nil.
Artist name and other data can be given directly or through an artist JSON file.
Example:
$ %s tx artist create-artist --artist="path/to/artist.json" --from mykey
Where artist.json contains:
{
  "name": "Freddy Mercury"
}
Which is equivalent to:
$ %s tx artist create-artist --name="Freddy Mercury" --from mykey
`,
			version.ClientName, version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			artist, err := parseCreateArtistFlags()
			if err != nil {
				return err
			}

			meta := types.MetaFromArtist(artist.Name)

			msg := types.NewMsgCreateArtist(meta, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagName, "", "Artist name")
	cmd.Flags().String(FlagArtist, "", "artist file path (if this path is given, other artist flags are ignored")

	return cmd
}

type artist struct {
	Name string
}

func parseCreateArtistFlags() (*artist, error) {
	artist := &artist{}
	artistFile := viper.GetString(FlagArtist)

	if artistFile == "" {
		artist.Name = viper.GetString(FlagName)
		return artist, nil
	}

	for _, flag := range ArtistFlags {
		if viper.GetString(flag) != "" {
			return nil, fmt.Errorf("--%s flag provided alongside --artist, which is a noop", flag)
		}
	}

	payload, err := ioutil.ReadFile(artistFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payload, artist)
	if err != nil {
		return nil, err
	}

	return artist, nil
}
