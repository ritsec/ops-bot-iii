package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
					Content: attendanceofMessage(u, span.Context()),
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

//returns the message sent to the user by the Attendanceof command
func attendanceofMessage(u *discordgo.User, ctx ddtrace.SpanContext) (message string) {
	span := tracer.StartSpan(
		"commands.slash.attendanceof:attendanceofMessage",
		tracer.ResourceName("/attendanceof:attendanceofMessage"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	message = ("**Signins for " + u.Username + "**")
	signinTypes := helpers.SigninTypeArray()

	totalSignins, err := data.Signin.GetSignins(u.ID, span.Context())
	if err != nil {
		logging.Error(nil, err.Error(), nil, span)
		totalSignins = 0
	}
	message += fmt.Sprintf("\n\tTotal Signins: `%d`", totalSignins)

	for _, signinType := range signinTypes {
		entSigninType := helpers.StringToType(signinType)
		signins, err := data.Signin.GetSigninsByType(u.ID, entSigninType, span.Context())
		if err != nil {
			logging.Error(nil, err.Error(), nil, span)
			signins = 0
		}
		if signins != 0 {
			message += fmt.Sprintf("\n\t%s: `%d`", signinType, signins)
		}
	}

	return message
}
