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
func attendanceMessage(userID string, ctx ddtrace.SpanContext) string {
	span := tracer.StartSpan(
		"commands.slash.attendance:attendanceMessage",
		tracer.ResourceName("/attendance:attendanceMessage"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	message := "**Your Signins**\n"
	signinTypes := map[string]int {
		"General Meeting": 0,
		"Contagion": 0,
		"IR": 0,
		"Ops": 0,
		"Ops IG": 0,
		"Red Team": 0,
		"Red Team Recruiting": 0,
		"RVAPT": 0,
		"Reversing": 0,
		"Physical": 0,
		"Wireless": 0,
		"WiCyS": 0,
		"Vulnerability Research": 0,
		"Mentorship": 0,
		"Other": 0
	}

	for signinType, signinCount := range signinTypes

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

}