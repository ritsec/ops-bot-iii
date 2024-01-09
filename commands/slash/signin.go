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

// Signin is a slash command that opens signins
func Signin() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
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
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			// Check if sign-in creator has already signed in
			recentSignin, err := data.Signin.RecentSignin(i.Member.User.ID, entSigninType, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			if !recentSignin {
				// Create sign-in for sign-in creator
				_, err = data.Signin.Create(i.Member.User.ID, entSigninType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}
			}

			if config.Google.Enabled {
				// Backup sign-in to google sheet
				err = google.SheetsAppendSignin(i.Member.User.ID, i.Member.User.Username, signinType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}
			}

			// Code Block that run when `Sign-in` button is pressed
			(*ComponentHandlers)[signinSlug] = func(s *discordgo.Session, j *discordgo.InteractionCreate) {
				span_signinSlug := tracer.StartSpan(
					"commands.slash.signin:Signin:signinSlug",
					tracer.ResourceName("/signin:signinSlug"),
					tracer.ChildOf(span.Context()),
				)
				defer span.Finish()

				// Check if user signed in recently
				recentSignin, err := data.Signin.RecentSignin(j.Member.User.ID, entSigninType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					return
				}

				if recentSignin {
					// User has already signed in, notify and exit
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

				// Create sign-in for user
				_, err = data.Signin.Create(j.Member.User.ID, entSigninType, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
					return
				}

				if config.Google.Enabled {
					// Backup sign-in to google sheet
					err = google.SheetsAppendSignin(j.Member.User.ID, j.Member.User.Username, signinType, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), j.Member.User, span_signinSlug, logrus.Fields{"error": err})
						return
					}
				}

				// Notify user that they have signed in
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

			// Send message with button to sign in
			message, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
				Content: "Signins are open for **" + signinType + "** until **" + time.Now().In(location).Add(time.Duration(delay)*time.Hour).Format("3:04PM") + "**!",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Signin",
								Style:    discordgo.SuccessButton,
								CustomID: signinSlug,
								Emoji: discordgo.ComponentEmoji{
									Name: "📝",
								},
							},
						},
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			// Notify sign-in creator that sign-in message was created
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Signin Message Created, it will close in %d hours!\n%s", delay, signinMessage(i.Member.User.ID, entSigninType, span.Context())),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			// Wait for sign-in to close
			time.Sleep(time.Duration(delay) * time.Hour)

			// Delete sign-in message
			err = s.ChannelMessageDelete(i.ChannelID, message.ID)
			if err != nil {
				logging.Error(s, "Error encounted while deleting message\n\n"+err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			// Get users who signed in
			userPairs, err := data.Signin.Query(time.Duration(12)*time.Hour, entSigninType, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			// Build message to send to sign-in creator
			msg := fmt.Sprintf("Signins for `%s`; %d users signed in:\n", signinType, len(userPairs))
			for _, user := range userPairs {
				msg += fmt.Sprintf("- %s\n", helpers.AtUser(user.Key))
			}

			if len(msg) <= 2000 {
				// Send full message to sign-in creator
				err = helpers.SendDirectMessageWithFile(s, i.Member.User.ID, msg, msg, span.Context())
			} else {
				// Send concatenated message to sign-in creator
				trimmedMsg := msg[:2000]
				trimmedMsg = trimmedMsg[:strings.LastIndex(trimmedMsg, "\n")]
				err = helpers.SendDirectMessageWithFile(s, i.Member.User.ID, trimmedMsg, msg, span.Context())
			}
			if err != nil {
				logging.Error(
					s, "Error encounted while sending direct message\n\n"+err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}
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

	// Get total signins and signins for type
	totalSignins, err := data.Signin.GetSignins(userID, span.Context())
	if err != nil {
		logging.Error(nil, err.Error(), nil, span)
		return fmt.Sprintf("You have sucessfully signed in for **%s**!", signinType)
	}

	// Get total signins and signins for type
	signins, err := data.Signin.GetSigninsByType(userID, signinType, span.Context())
	if err != nil {
		logging.Error(nil, err.Error(), nil, span)
		return fmt.Sprintf("You have sucessfully signed in for **%s**!\nYou have:\n\tTotal Signins: %d", signinType, totalSignins)
	}

	// Return full message
	return fmt.Sprintf("You have sucessfully signed in for **%s**!\nYou have:\n\tTotal Signins: `%d`\n\t%s Signins: `%d`", signinType, totalSignins, signinType, signins)
}
