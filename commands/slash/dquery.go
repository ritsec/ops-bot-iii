package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/ent/signin"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DQuery() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "dquery",
			Description:              "Query users by signins on specific date and outputs to CSV file",
			DefaultMemberPermissions: &permission.IGLead,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "type",
					Description: "The type of signin",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "All",
							Value: "All",
						},
						{
							Name:  "General Meeting",
							Value: "General Meeting",
						},
						{
							Name:  "Contagion",
							Value: "Contagion",
						},
						{
							Name:  "DFIR",
							Value: "DFIR",
						},
						{
							Name:  "Ops",
							Value: "Ops",
						},
						{
							Name:  "Ops IG",
							Value: "Ops IG",
						},
						{
							Name:  "Red Team",
							Value: "Red Team",
						},
						{
							Name:  "Red Team Recruiting",
							Value: "Red Team Recruiting",
						},
						{
							Name:  "Reversing",
							Value: "Reversing",
						},
						{
							Name:  "RVAPT",
							Value: "RVAPT",
						},
						{
							Name:  "Physical",
							Value: "Physical",
						},
						{
							Name:  "Vulnerability Research",
							Value: "Vulnerability Research",
						},
						{
							Name:  "Wireless",
							Value: "Wireless",
						},
						{
							Name:  "WiCyS",
							Value: "WiCyS",
						},
						{
							Name:  "Other",
							Value: "Other",
						},
					},
				},
				{
					Name:        "date",
					Description: "Specific date (YYYY-MM-DD) to query for",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.signin:Signin",
				tracer.ResourceName("/signin"),
			)
			defer span.Finish()

			signinType := i.ApplicationCommandData().Options[0].StringValue()
			dateRequested := i.ApplicationCommandData().Options[1].StringValue()

			var entSigninType signin.Type
			switch signinType {
			case "General Meeting":
				entSigninType = signin.TypeGeneralMeeting
			case "Contagion":
				entSigninType = signin.TypeContagion
			case "DFIR":
				entSigninType = signin.TypeDFIR
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
			case "Other":
				entSigninType = signin.TypeOther
			case "All":
				entSigninType = "All"
			}

			// Parsing the date as time.Time
			dateToQuery, err := time.Parse("2006-01-02", dateRequested)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			signins, err := data.Signin.DQuery(
				dateToQuery,
				entSigninType,
				span.Context(),
			)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			message := ""
			// Initial message
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "OBIII is processing...",
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			// Processing
			for x, signin := range signins {

				// Wait for 1 seconds after every 8 user's username is called
				if x > 0 && x%8 == 0 {
					time.Sleep(1 * time.Second)
				}

				user, err := s.User(signin.Key)

				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}

				if x == 0 {
					message += user.Username
				} else {
					message += fmt.Sprintf(",%s", user.Username)
				}
			}

			// Followup message
			followUpMessage := "Done"
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &followUpMessage,
				Files: []*discordgo.File{
					{
						Name:        "query.csv",
						ContentType: "text/csv",
						Reader:      strings.NewReader(message),
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}
		}
}
