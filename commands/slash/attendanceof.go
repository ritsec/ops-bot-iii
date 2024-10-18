package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
)

//Attendanceof slash command
func Attendanceof() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand {
		Name: 		 			  "attendanceof",
		Description: 			  "Get signin history of a user",
		DefaultMemberPermissions: &permission.IGLead,
		Options: []*discordgo.ApplicationCommandOption {
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user whose signin history you want to check",
				Required:    true,
			},
		},
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		span := tracer.StartSpan(
			"commands.slash.attendanceof:Attendanceof",
			tracer.ResourceName("/attendanceof"),
		)
		defer span.Finish()

		u := i.ApplicationCommandData().Options[0].UserValue(s)

		logging.Debug(s, "Attendanceof command received for " + u.Username, i.Member.User, span)

		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: attendanceMessage(u.ID, span.Context()),
					Flags: discordgo.MessageFlagsEphemeral,
				},
			},
		)
		if err != nil {
			logging.Error(s, err.Error(), i.Member.User, span)
		} else {
			logging.Debug(s, "Signin History Given for " + u.Username, i.Member.User, span)
		}
	}
}
