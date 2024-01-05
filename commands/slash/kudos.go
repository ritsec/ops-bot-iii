package slash

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// KudosChannelID is the channel ID to send kudos to
	KudosChannelID string = config.GetString("commands.kudos.channel_id")

	// KudosApprovalChannelID is the channel ID to approve kudos in
	KudosApprovalChannelID string = config.GetString("commands.kudos.approval_channel_id")
)

// Kudos is the kudos command
func Kudos() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "kudos",
			Description: "Send Kudos to a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The message to send about the user",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to send Kudos to",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "second_user",
					Description: "The second user to send Kudos to",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "third_user",
					Description: "The third user to send Kudos to",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "fourth_user",
					Description: "The fourth user to send Kudos to",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "fifth_user",
					Description: "The fifth user to send Kudos to",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "sixth_user",
					Description: "The sixth user to send Kudos to",
					Required:    false,
				},
			},
			DefaultMemberPermissions: &permission.Member,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.kudos:Kudos",
				tracer.ResourceName("/kudos"),
			)
			defer span.Finish()

			logging.Debug(s, "Kudos command received", i.Member.User, span)

			message := i.ApplicationCommandData().Options[0].StringValue()
			users := func() []*discordgo.User {
				var u []*discordgo.User
				for _, v := range i.ApplicationCommandData().Options[1:] {
					u = append(u, v.UserValue(s))
				}
				return u
			}()

			users_string := func() string {
				var u string

				if len(users) == 1 {
					u = users[0].Username
				} else if len(users) == 2 {
					u = users[0].Username + " and " + users[1].Username
				} else {
					for i, v := range users {
						if i == len(users)-1 {
							u += "and " + v.Username
						} else {
							u += v.Username + ", "
						}
					}
				}
				return u
			}()

			users_at := func() string {
				var u string
				for _, v := range users {
					u += helpers.AtUser(v.ID) + " "
				}
				return u
			}()

			logging.Debug(s, "Kudos sent to "+users_string+" for \""+message+"\"", i.Member.User, span)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Kudoes sent to " + users_string + " for \"" + message + "\"",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			author := &discordgo.MessageEmbedAuthor{}
			author.IconURL = i.Member.User.AvatarURL("")
			author.Name = i.Member.User.Username

			approve_slug := uuid.New().String()
			deny_slug := uuid.New().String()

			var wg sync.WaitGroup

			wg.Add(1)

			(*ComponentHandlers)[approve_slug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				span_approveSlug := tracer.StartSpan(
					"commands.slash.kudos:Kudos:approve_slug",
					tracer.ResourceName("/kudos:approve_slug"),
					tracer.ChildOf(span.Context()),
				)
				defer span_approveSlug.Finish()

				_, err := s.ChannelMessageSendComplex(KudosChannelID, &discordgo.MessageSend{
					Content: users_at,
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Kudos",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Message",
									Value: message,
								},
							},
						},
					},
				})

				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span_approveSlug, logrus.Fields{"error": err})
				} else {
					logging.Debug(s, "Kudos approved for "+users_string+" for \""+message+"\"", i.Member.User, span_approveSlug)
				}

				wg.Done()
			}

			(*ComponentHandlers)[deny_slug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				span_denySlug := tracer.StartSpan(
					"commands.slash.kudos:Kudos:deny_slug",
					tracer.ResourceName("/kudos:deny_slug"),
					tracer.ChildOf(span.Context()),
				)
				defer span_denySlug.Finish()

				logging.Debug(s, "Kudos denied for "+users_string+" for \""+message+"\"", i.Member.User, span_denySlug)

				wg.Done()
			}

			approvalMessage, err := s.ChannelMessageSendComplex(KudosApprovalChannelID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Author: author,
						Title:  "Kudos",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Message",
								Value:  message,
								Inline: false,
							},
						},
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Approve",
								Style:    discordgo.SuccessButton,
								CustomID: approve_slug,
								Emoji: discordgo.ComponentEmoji{
									Name: "✔️",
								},
							},
							discordgo.Button{
								Label:    "Deny",
								Style:    discordgo.DangerButton,
								CustomID: deny_slug,
								Emoji: discordgo.ComponentEmoji{
									Name: "✖️",
								},
							},
						},
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			wg.Wait()

			err = s.ChannelMessageDelete(KudosApprovalChannelID, approvalMessage.ID)
			if err != nil {
				logging.Error(s, "Error deleting channel message", i.Member.User, span, logrus.Fields{"error": err})
			}
			delete(*ComponentHandlers, approve_slug)
			delete(*ComponentHandlers, deny_slug)
		}
}
