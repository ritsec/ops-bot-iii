package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	GuildID string = config.GetString("guild_id")
	AlumniRequestChannelID string = config.GetString("commands.alumni-request.channel_id")
	MemberRoleID string = config.GetString("commands.member.member_role_id")
	AlumniRoleID string = config.GetString("commands.member.alumni_role_id")
)

// Alumni slash command
func Alumni() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "alumni",
			Description:              "Send request to E-Board to convert member role to alumni",
			DefaultMemberPermissions: &permission.Member,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.alumni:Alumni",
				tracer.ResourceName("/alumni"),
			)
			defer span.Finish()

			logging.Debug(s, "Alumni command received", i.Member.User, span)

			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Alumni Role Request Sent",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				},
			)
			if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				}

			approveSlug := uuid.New().String()
			denySlug := uuid.New().String()


			author := &discordgo.MessageEmbedAuthor{
				Name: i.Member.User.Username,
				IconURL: i.Member.User.Avatar,
			}

			//Send approval request message in e-board request channel
			message, err := s.ChannelMessageSendComplex(AlumniRequestChannelID, &discordgo.MessageSend{
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Approve",
								Style:    discordgo.SuccessButton,
								CustomID: approveSlug,
							},
							discordgo.Button{
								Label:    "Deny",
								Style:    discordgo.DangerButton,
								CustomID: denySlug,
							},
						},
					},
				},
				Embed: &discordgo.MessageEmbed{
					Description: "User is requesting to switch from role `Member` to role `Alumni`",
					Color: 0x95a5a6,
					Author: author,
				},
			})

			//Approve button handler
			(*ComponentHandlers)[approveSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				span_approveSlug := tracer.StartSpan(
					"commands.slash.alumni:Alumni:approveSlug",
					tracer.ResourceName("/alumni:approveSlug"),
					tracer.ChildOf(span.Context()),
				)
				defer span_approveSlug.Finish()

				err := s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, AlumniRoleID)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				}
				err = s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, MemberRoleID)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				}

				if err == nil {
					s.ChannelMessageEditComplex(&discordgo.MessageEdit{
						ID: message.ID,
						Channel: AlumniRequestChannelID,
						Embed: &discordgo.MessageEmbed{
							Description: "Alumni role switch approved and applied",
							Color: 0x57F287,
							Author: author,
						},
					})
				}
			}
			defer delete(*ComponentHandlers, approveSlug)

			//Deny button handler
			(*ComponentHandlers)[denySlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				span_denySlug := tracer.StartSpan(
					"commands.slash.alumni:Alumni:denySlug",
					tracer.ResourceName("/alumni:denySlug"),
					tracer.ChildOf(span.Context()),
				)
				defer span_denySlug.Finish()

				s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID: message.ID,
					Channel: AlumniRequestChannelID,
					Embed: &discordgo.MessageEmbed{
						Description: "Alumni role switch denied",
						Color: 0xED4245,
						Author: author,
					},
				})
			}
			defer delete(*ComponentHandlers, denySlug)

			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			} else {
				logging.Debug(s, "Alumni request sent", i.Member.User, span)
			}
		}
}
