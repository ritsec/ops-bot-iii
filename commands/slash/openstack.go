package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Openstack() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "openstack self service",
			Description:              "Create or reset your openstack account",
			DefaultMemberPermissions: &permission.Member,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option",
					Description: "Option of create or reset",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Create",
							Value: "Create",
						},
						{
							Name:  "Reset",
							Value: "Reset",
						},
					},
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.openstack:Openstack",
				tracer.ResourceName("/openstack"),
			)
			defer span.Finish()

			ssOption := i.ApplicationCommandData().Options[0].StringValue()

			// CHECK IF USER IS DM'ABLE
			err := helpers.SendDirectMessage(s, i.Member.User.ID, "Checking to see if your DMs are open... your openstack account password will be sent here!", span.Context())
			if err != nil {
				logging.Debug(s, "User's DMs are not open", i.Member.User, span)
				err = s.InteractionRespond(
					i.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Your DMs are not open! Please open your DMs and run the command again",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					},
				)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
			}

			// CHECK IF USER EXISTS ON OPENSTACK ALREADY
			err, result = helpers.CheckIfExists()

			if ssOption == "Create" {
				// CREATE THE ACCOUNT
			} else if ssOption == "Reset" {
				// RESET THE PASSWORD OF THE ACCOUNT
			} else {
				// GRACEFULLY CLOSE
				logging.Error(s, "User somehow typed a different option for /openstack?", i.Member.User, span)
			}
		}
}
