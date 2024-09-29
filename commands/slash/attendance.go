package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/ent/signin"
	"github.com/ritsec/ops-bot-iii/google"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

//Attendance slash command
func Attendance() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand {
		Name: 		 			  "attendance",
		Description: 			  "Get signin history",
		DefaultMemberPermissions: &permission.Member,
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		span := tracer.StartSpan(
			"commands.slash.attendance:Attendance",
			tracer.ResourceName("/attendance"),
		)
		defer span.Finish()

		logging.Debug(s, "Attendance command received", i.Member.User, span)

		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: (

					),
					Flags: discordgo.MessageFlagsEphemeral,
				},
			},
		)
		if err != nil {
			logging.Error(s, err.Error(), i.Member.User, span)
		} else {
			logging.Debug(s, "Signin History Given", i.Member.User, span)
		}
	}
}

//returns the message sent to the user by the Attendance command
func attendanceMessage(userID string, ctx ddtrace.SpanContext) (message string) {
	span := tracer.StartSpan(
		"commands.slash.attendance:attendanceMessage",
		tracer.ResourceName("/attendance:attendanceMessage"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	message := "**Your Signins:**"
	signinTypes := [...]string{
		"General Meeting",
		"Contagion",
		"IR",
		"Ops",
		"Ops IG",
		"Red Team",
		"Red Team Recruiting",
		"RVAPT",
		"Reversing",
		"Physical",
		"Wireless",
		"WiCyS",
		"Vulnerability Research",
		"Mentorship",
		"Other",
	}

	totalSignins, err := data.Signin.GetSignins(userID, span.Context())
	if err != nil {
		logging.Error(nil, err.Error(), nil, span)
		totalSignins = 0
	}
	message += fmt.Sprintf("\n\tTotal Signins: `%d`", totalSignins)

	for _, signinType := range signinTypes {
		var entSigninType signin.Type
			switch signinType {
			case "General Meeting":
				entSigninType = signin.TypeGeneralMeeting
			case "Contagion":
				entSigninType = signin.TypeContagion
			case "IR":
				entSigninType = signin.TypeIR
			case "Ops":
				entSigninType = signin.TypeOps
			case "Ops IG":
				entSigninType = signin.TypeOpsIG
			case "Red Team":
				entSigninType = signin.TypeRedTeam
			case "Red Team Recruiting":
				entSigninType = signin.TypeRedTeamRecruiting
			case "RVAPT":
				entSigninType = signin.TypeRVAPT
			case "Reversing":
				entSigninType = signin.TypeReversing
			case "Physical":
				entSigninType = signin.TypePhysical
			case "Wireless":
				entSigninType = signin.TypeWireless
			case "WiCyS":
				entSigninType = signin.TypeWiCyS
			case "Vulnerability Research":
				entSigninType = signin.TypeVulnerabilityResearch
			case "Mentorship":
				entSigninType = signin.TypeMentorship
			case "Other":
				entSigninType = signin.TypeOther
			}
		signins, err := data.Signin.GetSigninsByType(userID, entSigninType, span.Context())
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