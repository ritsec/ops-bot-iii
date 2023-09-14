package slash

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/slash/permission"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/data"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent/signin"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/google"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/helpers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Signin is a slash command that opens signins
func Signin() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "signin",
			Description: "Open Signins",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "type",
					Description: "The type of signin",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
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

			logging.Debug(s, fmt.Sprintf("Signin Created: %s", signinType), i.Member.User, span)

			location, err := time.LoadLocation("America/New_York")
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			signinSlug := uuid.New().String()

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
			}

			(*ComponentHandlers)[signinSlug] = func(s *discordgo.Session, j *discordgo.InteractionCreate) {
				span_signinSlug := tracer.StartSpan(
					"commands.slash.signin:Signin:signinSlug",
					tracer.ResourceName("/signin:signinSlug"),
					tracer.ChildOf(span.Context()),
				)
				defer span.Finish()

				recentSignin, err := data.Signin.RecentSignin(j.Member.User.ID, entSigninType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					return
				}

				if recentSignin {
					err = s.InteractionRespond(j.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You have already signed in for **" + signinType + "**!",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					}
					return
				}

				_, err = data.Signin.Create(j.Member.User.ID, entSigninType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					return
				}

				err = google.SheetsAppendSignin(j.Member.User.ID, j.Member.User.Username, signinType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					return
				}

				err = s.InteractionRespond(j.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: signinMessage(j.Member.User.ID, entSigninType, span_signinSlug.Context()),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					return
				}
			}
			defer delete(*ComponentHandlers, signinSlug)

			var delay int
			if signinType == "General Meeting" {
				delay = 4
			} else {
				delay = 2
			}

			message, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
				Content: "Signins are open for **" + signinType + "** until **" + time.Now().In(location).Add(time.Duration(delay)*time.Hour).Format("3:04PM") + "**!",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Signin",
								Style:    discordgo.SuccessButton,
								CustomID: signinSlug,
							},
						},
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Signin Message Created, it will close in %d hours!", delay),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			defer func() {
				err = s.ChannelMessageDelete(i.ChannelID, message.ID)
				if err != nil {
					logging.Error(s, "Error encounted while deleting message\n\n"+err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				}

				users, err := data.Signin.QueryUsers(time.Duration(12)*time.Hour, entSigninType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}

				message := fmt.Sprintf("Signins for `%s`; %d users signed in:\n", signinType, len(users))
				for _, user := range users {
					message += fmt.Sprintf("- %s\n", helpers.AtUser(user.ID))
				}

				err = helpers.SendDirectMessage(s, i.Message.Author.ID, "", span.Context())
				if err != nil {
					logging.Error(s, "Error encounted while sending direct message\n\n"+err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				}
			}()

			time.Sleep(time.Duration(delay) * time.Hour)
		},
	}
}

// signinMessage returns the message to send to the user after they have signed in
func signinMessage(userID string, signinType signin.Type, ctx ddtrace.SpanContext) string {
	span := tracer.StartSpan(
		"commands.slash.signin:signinMessage",
		tracer.ResourceName("/signin:signinMessage"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	totalSignins, err := data.Signin.GetSignins(userID, span.Context())
	if err != nil {
		logging.Error(nil, err.Error(), nil, span)
		return fmt.Sprintf("You have sucessfully signed in for **%s**!", signinType)
	}

	signins, err := data.Signin.GetSigninsByType(userID, signinType, span.Context())
	if err != nil {
		logging.Error(nil, err.Error(), nil, span)
		return fmt.Sprintf("You have sucessfully signed in for **%s**!\nYou have:\n\tTotal Signins: %d", signinType, totalSignins)
	}

	return fmt.Sprintf("You have sucessfully signed in for **%s**!\nYou have:\n\tTotal Signins: `%d`\n\t%s Signins: `%d`", signinType, totalSignins, signinType, signins)
}
