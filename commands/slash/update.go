package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Update() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "update",
			Description:              "Update the bot",
			DefaultMemberPermissions: &permission.Admin,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "force",
					Description: "Force the bot to reboot",
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Required:    false,
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.update:Update",
				tracer.ResourceName("/update"),
			)
			defer span.Finish()

			logging.Debug(s, "Update command received", i.Member.User, span)

			force := false
			if len(i.ApplicationCommandData().Options) != 0 {
				force = i.ApplicationCommandData().Options[0].BoolValue()
			}

			update, err := helpers.UpdateMainBranch()
			if err != nil {
				logging.Error(s, "Error updating main branch", i.Member.User, span, logrus.Fields{"err": err.Error()})
				return
			}

			if !update {
				logging.Debug(s, "No update available", i.Member.User, span)

				if !force {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No update available\nIf you want to force an update, use `/update force`",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						logging.Error(s, "Error responding to interaction", i.Member.User, span, logrus.Fields{"err": err.Error()})
						return
					}

					return
				} else {
					logging.Debug(s, "Forcing update", i.Member.User, span)

					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No update available; forcing update\nBot will be up temporarily once done updating",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						logging.Error(s, "Error responding to interaction", i.Member.User, span, logrus.Fields{"err": err.Error()})
						return
					}

					err = helpers.BuildOBIII()
					if err != nil {
						logging.Error(s, "Error building OBIII", i.Member.User, span, logrus.Fields{"err": err.Error()})

						content := fmt.Sprintf("Error building OBIII\n\nError:\n%s", err.Error())
						_, err = s.InteractionResponseEdit(
							i.Interaction,
							&discordgo.WebhookEdit{
								Content: &content,
							},
						)
						if err != nil {
							logging.Error(s, "Error editing interaction response", i.Member.User, span, logrus.Fields{"err": err.Error()})
						}

						return
					}

					err = helpers.Exit()
					if err != nil {
						logging.Error(s, "Error exiting", i.Member.User, span, logrus.Fields{"err": err.Error()})
						return
					}

					return
				}
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Update available; restarting now\nBot will be up temporarily once done updating",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, "Error responding to interaction", i.Member.User, span, logrus.Fields{"err": err.Error()})
				return
			}

			err = helpers.BuildOBIII()
			if err != nil {
				content := fmt.Sprintf("Error building OBIII\n\nError:\n%s", err.Error())
				_, err = s.InteractionResponseEdit(
					i.Interaction,
					&discordgo.WebhookEdit{
						Content: &content,
					},
				)

				logging.Error(s, "Error building OBIII", i.Member.User, span, logrus.Fields{"err": err.Error()})
				return
			}

			err = helpers.Exit()
			if err != nil {
				logging.Error(s, "Error exiting", i.Member.User, span, logrus.Fields{"err": err.Error()})
				return
			}
		}
}
