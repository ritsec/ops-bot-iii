package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Reboot is a slash command that reboots the bot
func Scoreboard() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "scoreboard",
			Description:              "Shitposting Scoreboard",
			DefaultMemberPermissions: &permission.Member,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.scoreboard:Scoreboard",
				tracer.ResourceName("/scoreboard"),
			)
			defer span.Finish()

			logging.Debug(s, "Scoreboard command received", i.Member.User, span)

			entShitposts, err := data.Shitposts.GetTopXShitposts(10, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
				return
			}

			posts := ""

			for i, entShitpost := range entShitposts {
				posts += fmt.Sprintf("%d. %s [link](%s) - %d\n", i+1, helpers.AtUser(entShitpost.Edges.User.ID), helpers.JumpURLByID(entShitpost.ChannelID, entShitpost.ID), entShitpost.Count)
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Top 10 Shitposts:\n" + posts,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
				return
			}
		}
}
