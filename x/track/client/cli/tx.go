package cli

import (
	"bufio"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
	"strings"
)

const (
	flagDao      = "dao"
	flagArtist   = "artist"
	flagDuration = "duration"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	contentTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	contentTxCmd.AddCommand(flags.PostCommands(
		GetCmdAdd(cdc),
	)...)

	return contentTxCmd
}

func GetCmdAdd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new track to bitsong",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a new track to bitsong.
Example:
$ %s tx track add [title] [duration] \
--dao "2000:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a" \
--dao "1000:bitsong13r9ryyfltaz8rsqqumqxusgtw0ne4udhxm5jm4" \
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			title := args[0]

			artistsStr, err := cmd.Flags().GetStringArray(flagArtist)
			if err != nil {
				return fmt.Errorf("invalid flag value: %s", flagArtist)
			}

			artists := make([]string, len(artistsStr))
			for i, artist := range artistsStr {
				artists[i] = artist
			}

			number := uint(1) // default 1

			duration, err := cmd.Flags().GetUint(flagDuration)
			if err != nil {
				return err
			}
			if duration < 15000 {
				return fmt.Errorf("duration must be > 15000 milliseconds")
			}

			explicit := false // default false
			pUrl := ""        // default empty

			daoStr, err := cmd.Flags().GetStringArray(flagDao)
			if err != nil {
				return fmt.Errorf("invalid dao value")
			}

			dao := types.Dao{}
			for _, rh := range daoStr {
				deArgs := strings.Split(rh, ":")
				if len(deArgs) != 2 {
					return fmt.Errorf("the dao format must be \"shares:address\" ex: \"1000:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a\"")
				}
				des, err := sdk.NewDecFromStr(deArgs[0])
				if err != nil {
					return err
				}
				deAddr, err := sdk.AccAddressFromBech32(deArgs[1])
				if err != nil {
					return fmt.Errorf("dao entity address is wrong, %s", err.Error())
				}
				de := types.NewDaoEntity(des, deAddr)
				dao = append(dao, de)
			}

			msg := types.NewMsgTrackAdd(title, artists, number, uint(duration), explicit, nil, nil, pUrl, dao)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringArray(flagArtist, []string{}, "Track Artists")
	cmd.Flags().StringArray(flagDao, []string{}, "Track DAO")
	cmd.Flags().Uint(flagDuration, 0, "Track duration in milliseconds")

	return cmd
}
