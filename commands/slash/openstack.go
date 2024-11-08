package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Openstack() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "openstack",
			Description:              "Create or reset your openstack account",
			DefaultMemberPermissions: &permission.Member,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option",
					Description: "Option of create or reset",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Create",
							Value: "Create",
						},
						{
							Name:  "Reset",
							Value: "Reset",
						},
					},
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.openstack:Openstack",
				tracer.ResourceName("/openstack"),
			)
			defer span.Finish()

			ssOption := i.ApplicationCommandData().Options[0].StringValue()

			// err := helpers.SourceOpenRC()
			err := helpers.DebugSourceOpenRC(s, i.Member.User, span)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}

			// CHECK IF USER IS DM'ABLE
			err = helpers.SendDirectMessage(s, i.Member.User.ID, "Checking to see if your DMs are open... your openstack account username and password will be sent here!", span.Context())
			if err != nil {
				logging.Debug(s, "User's DMs are not open", i.Member.User, span)
				err = s.InteractionRespond(
					i.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Your DMs are not open! Please open your DMs and run the command again.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					},
				)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				return
			}

			// GET EMAIL AND CHECK IF IT IS VALID
			email, err := data.User.GetEmail(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}
			if email == "" {
				logging.Debug(s, "User has no email", i.Member.User, span)
				err = s.InteractionRespond(
					i.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You have no verified email. Run /member and verify your email and run this command again.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					},
				)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				return
			}

			// CHECK IF USER EXISTS ON OPENSTACK ALREADY
			exists, err := helpers.DebugCheckIfExists(s, i.Member.User, span, email)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
				return
			}

			if ssOption == "Create" {
				if exists {
					logging.Debug(s, "User already has an openstack account and is trying to create one", i.Member.User, span)
					err = s.InteractionRespond(
						i.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Openstack account already exisits. Run the reset option if you forgot your password.",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						},
					)
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					return
				}

				// CREATE THE ACCOUNT
				// username, password, err := helpers.Create(email)
				username, password, err := helpers.DebugCreate(s, i.Member.User, span, email)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}

				// SEND THE USERNAME AND PASSWORD TO THE USER VIA DM
				message := fmt.Sprintf("Thank you for reaching out to us!\nHere are your credentials for RITSEC's Openstack:\n\nUsername: %s\nTemporary Password: %s\n\nPlease change the password\nOpenstack link: stack.ritsec.cloud", username, password)
				err = helpers.SendDirectMessage(s, i.Member.User.ID, message, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}
			} else if ssOption == "Reset" {
				// CHECK IF USER IS TRYING TO RESET PASSWORD ON NON-EXISTENT ACCOUNT
				if !exists {
					logging.Debug(s, "User does not have an openstack account and is trying to reset the password on it", i.Member.User, span)
					err = s.InteractionRespond(
						i.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Openstack account does not exist and you are trying to reset it.",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						},
					)
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					return
				}

				// RESET THE PASSWORD OF THE ACCOUNT
				username, password, err := helpers.Reset(email)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}

				message := fmt.Sprintf("Thank you for reaching out to us!\nHere are your credentials for RITSEC's Openstack:\n\nUsername: %s\nTemporary Password: %s\n\nPlease change the password\nOpenstack link: stack.ritsec.cloud", username, password)
				err = helpers.SendDirectMessage(s, i.Member.User.ID, message, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}
			} else {
				return
			}
		}
}
