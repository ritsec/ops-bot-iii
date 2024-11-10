package slash

import (
	"fmt"
	"strings"

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
			helpers.InitialMessage(s, i, fmt.Sprintf("You ran the /openstack command to %s your account!", strings.ToLower(ssOption)))

			// Initialize the environment variables for Openstack CLI
			helpers.SetOpenstackRC()

			err := helpers.UpdateMessage(s, i, "Checking if your DMs are open...")
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}

			// Checking if the user is DM'able
			err = helpers.SendDirectMessage(s, i.Member.User.ID, "Checking to see if your DMs are open... your openstack account username and password will be sent here!", span.Context())
			if err != nil {
				logging.Debug(s, "User's DMs are not open", i.Member.User, span)
				err = helpers.UpdateMessage(s, i, "Your DMs are not open! Please open your DMs and run the command again.")
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				return
			}

			// Get email and check if it is an actual email
			email, err := data.User.GetEmail(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}
			if email == "" {
				logging.Debug(s, "User has no email", i.Member.User, span)
				err = helpers.UpdateMessage(s, i, "You have no verified email. Run /member and verify your email and run this command again.")
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				return
			}

			// Check if user exists on Openstack already
			exists, err := helpers.CheckIfExists(email)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
				return
			}

			if ssOption == "Create" {
				// Check if user trying to create an account when it already has one
				if exists {
					logging.Debug(s, "User already has an openstack account and is trying to create one", i.Member.User, span)
					err = helpers.UpdateMessage(s, i, "Openstack account already exisits. Run the reset option if you forgot your password.")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					return
				}

				// Create the account
				username, password, err := helpers.Create(email)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}

				// Send the username and password to the usuer via DM
				message := fmt.Sprintf("Thank you for reaching out to us!\nHere are your credentials for RITSEC's Openstack:\n\nUsername: %s\nTemporary Password: %s\n\nPlease change the password\nOpenstack link: stack.ritsec.cloud", username, password)
				logging.Debug(s, "Sent username and password to member", i.Member.User, span)
				err = helpers.SendDirectMessage(s, i.Member.User.ID, message, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}
			} else if ssOption == "Reset" {
				// Check if the user is trying to reset password on non-existent account
				if !exists {
					logging.Debug(s, "User does not have an openstack account and is trying to reset the password on it", i.Member.User, span)
					err = helpers.UpdateMessage(s, i, "Openstack account does not exist and you are trying to reset it.")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					return
				}

				// Reset the password of the account
				username, password, err := helpers.Reset(email)
				logging.Debug(s, "User has the openstack account password reset", i.Member.User, span)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}

				message := fmt.Sprintf("Thank you for reaching out to us!\n Here are your credentials for RITSEC's Openstack:\n\nUsername: %s\nTemporary Password: %s\n\nPlease change the password\nOpenstack link: stack.ritsec.cloud", username, password)
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
