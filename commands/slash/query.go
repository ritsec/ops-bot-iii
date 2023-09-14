package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/slash/permission"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/data"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent/signin"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/helpers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Query users based on signins
func Query() *structs.SlashCommand {
	minValue := float64(0)
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "query",
			Description: "Query users by signins",
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
					Name:        "hours",
					Description: "The number of hours to query for",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
					MinValue:    &minValue,
				},
				{
					Name:        "days",
					Description: "The number of days to query for",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
					MinValue:    &minValue,
				},
				{
					Name:        "weeks",
					Description: "The number of weeks to query for",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
					MinValue:    &minValue,
				},
			},
			DefaultMemberPermissions: &permission.IGLead,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.signin:Signin",
				tracer.ResourceName("/signin"),
			)
			defer span.Finish()

			signinType := i.ApplicationCommandData().Options[0].StringValue()

			var (
				hours int
				days  int
				weeks int
			)

			if len(i.ApplicationCommandData().Options) > 1 {
				for _, option := range i.ApplicationCommandData().Options[1:] {
					switch option.Name {
					case "hours":
						hours = int(option.IntValue())
					case "days":
						days = int(option.IntValue())
					case "weeks":
						weeks = int(option.IntValue())
					}
				}
			}

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

			signins, err := data.Signin.Query(
				time.Duration(hours)*time.Hour+time.Duration(days)*24*time.Hour+time.Duration(weeks)*7*24*time.Hour,
				entSigninType,
				span.Context(),
			)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			sum := 0
			for _, signin := range signins {
				sum += signin.Value
			}

			message := fmt.Sprintf("Signin Type: `%s`\nTotal Signins: `%d`\nTime Delta: `hours=%d,days=%d,weeks=%d`\n", signinType, sum, hours, days, weeks)

			for _, signin := range signins {
				username, err := helpers.Username(s, signin.Key)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					username = "Failed to resolve"
				}
				message += fmt.Sprintf("[%d] [%s] %s\n", signin.Value, username, helpers.AtUser(signin.Key))
			}
			if len(message) <= 2000 {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: message,
						Flags:   discordgo.MessageFlagsEphemeral,
						Files: []*discordgo.File{
							{
								Name:        "query.txt",
								ContentType: "text/plain",
								Reader:      strings.NewReader(message),
							},
						},
					},
				})
			} else {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Files: []*discordgo.File{
							{
								Name:        "query.txt",
								ContentType: "text/plain",
								Reader:      strings.NewReader(message),
							},
						},
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
			}
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}
		},
	}
}
