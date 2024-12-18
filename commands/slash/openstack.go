package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/osclient"
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
			if config.Openstack.Enabled {
				span := tracer.StartSpan(
					"commands.slash.openstack:Openstack",
					tracer.ResourceName("/openstack"),
				)
				defer span.Finish()

				ssOption := i.ApplicationCommandData().Options[0].StringValue()
				err := helpers.InitialMessage(s, i, fmt.Sprintf("You ran the /openstack command to %s your account!", strings.ToLower(ssOption)))
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}

				err = helpers.UpdateMessage(s, i, "Checking your email...")
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				// Get email and check if it is an actual email
				email, err := data.User.GetEmail(i.Member.User.ID, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				if email == "" {
					logging.Debug(s, "User has no email", i.Member.User, span)
					err = helpers.UpdateMessage(s, i, "You have no verified email. Run /verify to add your email and run this command again.")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					return
				}

				// Check if user exists on Openstack already
				exists, err := osclient.CheckUserExists(email)
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
					return
				}

				if ssOption == "Create" {
					// Stop here if user trying to create an account when it already has one
					if exists {
						logging.Debug(s, "User already has an openstack account and is trying to create one", i.Member.User, span)
						err = helpers.UpdateMessage(s, i, "Openstack account already exisits. Run the reset option if you forgot your password.")
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span)
						}
						return
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

					// Create the account
					err = helpers.UpdateMessage(s, i, "Creating your account...")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					username, password, err := osclient.Create(email)
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

					err = helpers.UpdateMessage(s, i, "Sent the username and password to your DMs, check your DMs!")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
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

					// Check to see if the user's reset timestamp exists
					ts_exists, err := data.Openstack.Exists(i.Member.User.ID, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
						return
					}
					if !ts_exists {
						// Create the row for user's timestamp and don't check it
						_, err = data.Openstack.Create(i.Member.User.ID, span.Context())
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span)
							return
						}
					} else {
						// Check the user's timestamp if has done recently

						// Get the user's openstacks info
						openstack_ent, err := data.Openstack.Get(i.Member.User.ID, span.Context())
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span)
							return
						}
						tx, err := time.LoadLocation("America/New_York")
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span)
							return
						}
						// Check if the user has done the reset recently
						if time.Now().In(tx).Sub(openstack_ent.Timestamp) >= -2*time.Hour && time.Now().In(tx).Sub(openstack_ent.Timestamp) <= 2*time.Hour {
							err = helpers.UpdateMessage(s, i, "You have tried to reset too much in the past 2 hours, please wait before trying again!")
							if err != nil {
								logging.Error(s, err.Error(), i.Member.User, span)
							}
							return
						}
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

					// 	// Reset the password of the account
					err = helpers.UpdateMessage(s, i, "Resetting your account...")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}
					username, password, err := osclient.Reset(email)
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

					err = helpers.UpdateMessage(s, i, "Sent the username and password to your DMs, check your DMs!")
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
					}

					// Update the timestamp for the last reset for the user
					_, err = data.Openstack.Update(i.Member.User.ID, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span)
						return
					}
				} else {
					return
				}
			}
		}
}
