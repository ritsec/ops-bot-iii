package slash

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/slash/permission"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/helpers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// KudosChannelID is the channel ID to send kudos to
	KudosChannelID string = config.GetString("commands.kudos.channel_id")

	// KudosApprovalChannelID is the channel ID to approve kudos in
	KudosApprovalChannelID string = config.GetString("commands.kudos.approval_channel_id")
)

// Kudos is the kudos command
func Kudos() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
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
			},
			DefaultMemberPermissions: &permission.Member,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.kudos:Kudos",
				tracer.ResourceName("/kudos"),
			)
			defer span.Finish()

			logging.Debug(s, "Kudos command received", i.Member.User, span)

			message := i.ApplicationCommandData().Options[0].StringValue()
			user := i.ApplicationCommandData().Options[1].UserValue(s)

			logging.Debug(s, "Kudos sent to "+user.Username+" for \""+message+"\"", i.Member.User, span)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Kudoes sent to " + user.Username + " for \"" + message + "\"",
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
					Content: helpers.AtUser(user.ID),
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
					logging.Debug(s, "Kudos approved for "+user.Username+" for \""+message+"\"", i.Member.User, span_approveSlug)
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

				logging.Debug(s, "Kudos denied for "+user.Username+" for \""+message+"\"", i.Member.User, span_denySlug)

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
							},
							discordgo.Button{
								Label:    "Deny",
								Style:    discordgo.DangerButton,
								CustomID: deny_slug,
							},
						},
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			wg.Wait()

			s.ChannelMessageDelete(KudosApprovalChannelID, approvalMessage.ID)
			delete(*ComponentHandlers, approve_slug)
			delete(*ComponentHandlers, deny_slug)
		},
	}
}
