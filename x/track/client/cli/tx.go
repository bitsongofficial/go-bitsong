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
	"github.com/spf13/viper"
	"strings"
)

const (
	flagDao        = "dao"
	flagArtist     = "artist"
	flagFeat       = "feat"
	flagProducer   = "producer"
	flagGenre      = "genre"
	flagMood       = "mood"
	flagTag        = "tag"
	flagExplicit   = "explicit"
	flagLabel      = "label"
	flagCredits    = "credits"
	flagCopyright  = "copyright"
	flagPreviewUrl = "preview-url"
	flagDuration   = "duration"
	flagNumber     = "number"
	flagExtID      = "ext-id"
	flagExtUrl     = "ext-url"
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
$ %s tx track add [title] \
--artist "Dj Angelo" \
--artist "Angelo 2" \
--feat "Singer 1" \
--feat "Singer 2" \
--producer "The best Producer" \
--producer "The Cat" \
--tag "tag 1" \
--tag "tag 2" \
--genre "Pop" \
--mood "Energetic" \
--label "My indie label" \
--credits "Thanks to..." \
--copyright "Creative Commons" \
--preview-url "https://my-preview-link" \
--duration 15001 \
--number 1 \
--dao "100:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a" \
--ext-id "youtube=M7Iwkxy_Hjw" --ext-id "spotify=12487825" \
--ext-url "youtube=https://www.youtube.com/watch?v=y6veLh2b1Js" --ext-url "soundcloud=https://soundcloud.com/bangtan/thankyouarmy2020"
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

			featStr, err := cmd.Flags().GetStringArray(flagFeat)
			if err != nil {
				return fmt.Errorf("invalid flag value: %s", flagFeat)
			}

			feats := make([]string, len(featStr))
			for i, feat := range featStr {
				feats[i] = feat
			}

			producersStr, err := cmd.Flags().GetStringArray(flagProducer)
			if err != nil {
				return fmt.Errorf("invalid flag value: %s", flagProducer)
			}

			producers := make([]string, len(producersStr))
			for i, producer := range producersStr {
				producers[i] = producer
			}

			genre, err := cmd.Flags().GetString(flagGenre)
			if err != nil {
				return fmt.Errorf("invalid flag value: %s", flagGenre)
			}

			mood, err := cmd.Flags().GetString(flagMood)
			if err != nil {
				return fmt.Errorf("invalid flag value: %s", flagMood)
			}

			tagsStr, err := cmd.Flags().GetStringArray(flagTag)
			if err != nil {
				return fmt.Errorf("invalid flag value: %s", flagTag)
			}

			tags := make([]string, len(tagsStr))
			for i, tag := range tagsStr {
				tags[i] = tag
			}

			number, err := cmd.Flags().GetUint(flagNumber)
			if err != nil {
				return err
			}

			duration, err := cmd.Flags().GetUint(flagDuration)
			if err != nil {
				return err
			}
			if duration < 15000 {
				return fmt.Errorf("duration must be > 15000 milliseconds")
			}

			explicit := viper.GetBool(flagExplicit)
			label := viper.GetString(flagLabel)
			credits := viper.GetString(flagCredits)
			copyright := viper.GetString(flagCopyright)
			pUrl := viper.GetString(flagPreviewUrl)

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

			extIdsStr, err := cmd.Flags().GetStringArray(flagExtID)
			extIds := make(map[string]string)
			for _, extID := range extIdsStr {
				idArgz := strings.SplitN(extID, "=", 2)
				if len(idArgz) != 2 {
					return fmt.Errorf("the ext format must be \"key=value\" e.g.: \"youtube=M7Iwkxy_Hjw\"")
				}
				extIds[idArgz[0]] = idArgz[1]
			}
			if len(extIds) == 0 {
				extIds = nil
			}

			extUrlsStr, err := cmd.Flags().GetStringArray(flagExtUrl)
			extUrls := make(map[string]string)
			for _, extUrl := range extUrlsStr {
				urlArgz := strings.SplitN(extUrl, "=", 2)
				if len(urlArgz) != 2 {
					return fmt.Errorf("the ext format must be \"key=value\" e.g.: \"youtube=https://www.youtube.com/watch?v=y6veLh2b1Js\"")
				}
				extUrls[urlArgz[0]] = urlArgz[1]
			}
			if len(extUrls) == 0 {
				extUrls = nil
			}

			msg := types.NewMsgTrackAdd(title, artists, feats, producers, tags, genre, mood, label, credits, copyright,
				pUrl, number, duration, explicit, extIds, extUrls, dao,
			)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringArray(flagArtist, []string{}, "Track Artists")
	cmd.Flags().StringArray(flagFeat, []string{}, "Track Feat")
	cmd.Flags().StringArray(flagProducer, []string{}, "Track Producers")
	cmd.Flags().String(flagGenre, "", "Track Genre")
	cmd.Flags().String(flagMood, "", "Track Mood")
	cmd.Flags().String(flagLabel, "", "Track Label")
	cmd.Flags().String(flagCredits, "", "Track Credits")
	cmd.Flags().String(flagCopyright, "", "Track Copyright")
	cmd.Flags().String(flagPreviewUrl, "", "Track Preview Url")
	cmd.Flags().StringArray(flagTag, []string{}, "Track Tags")
	cmd.Flags().Bool(flagExplicit, false, "Track explicit (true | false)")
	cmd.Flags().StringArray(flagDao, []string{}, "Track DAO")
	cmd.Flags().StringArray(flagExtID, []string{}, "Track External ids")
	cmd.Flags().StringArray(flagExtUrl, []string{}, "Track External URLs")
	cmd.Flags().Uint(flagDuration, 0, "Track duration in milliseconds")
	cmd.Flags().Uint(flagNumber, 1, "Track Number (default is 1)")

	return cmd
}
