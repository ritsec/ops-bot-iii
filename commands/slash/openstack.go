package slash

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// Environment variables required for openstack CLI
	OS_AUTH_URL             string = config.GetString("openstack.ENV.OS_AUTH_URL")
	OS_PROJECT_ID           string = config.GetString("openstack.ENV.OS_PROJECT_ID")
	OS_PROJECT_NAME         string = config.GetString("openstack.ENV.OS_PROJECT_NAME")
	OS_USER_DOMAIN_NAME     string = config.GetString("openstack.ENV.OS_USER_DOMAIN_NAME")
	OS_PROJECT_DOMAIN_ID    string = config.GetString("openstack.ENV.OS_PROJECT_DOMAIN_ID")
	OS_USERNAME             string = config.GetString("openstack.ENV.OS_USERNAME")
	OS_PASSWORD             string = config.GetString("openstack.ENV.OS_PASSWORD")
	OS_REGION_NAME          string = config.GetString("openstack.ENV.OS_REGION_NAME")
	OS_INTERFACE            string = config.GetString("openstack.ENV.OS_INTERFACE")
	OS_IDENTITY_API_VERSION string = config.GetString("openstack.ENV.OS_IDENTITY_API_VERSION")

	// Paths for scripts to automate openstack user management
	new_member      string = config.GetString("openstack.SCRIPTS.new_member")
	reset_password  string = config.GetString("openstack.SCRIPTS.reset_password")
	check_if_exists string = config.GetString("openstack.SCRIPTS.check_if_exists")
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
			err := helpers.InitialMessage(s, i, fmt.Sprintf("You ran the /openstack command to %s your account!", strings.ToLower(ssOption)))
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}

			// Initialize the environment variables for Openstack CLI
			SetOpenstackRC()

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
				err = helpers.UpdateMessage(s, i, "You have no verified email. Run /member and verify your email and run this command again.")
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				return
			}

			// Check if user exists on Openstack already
			exists, err := CheckIfExists(email)
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
				username, password, err := Create(email)
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

				// Reset the password of the account
				err = helpers.UpdateMessage(s, i, "Resetting your account...")
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span)
				}
				username, password, err := Reset(email)
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
			} else {
				return
			}
		}
}

func Create(email string) (username string, password string, error error) {
	createCmd := exec.Command(new_member, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCmd.Stdout = stdout
	createCmd.Stderr = stderr

	err := createCmd.Start()
	if err != nil {
		return "", "", err
	}
	err = createCmd.Wait()
	if err != nil {
		return "", "", err
	}

	output := strings.Fields(stdout.String())
	username = output[0]
	password = output[1]

	return username, password, nil
}

func Reset(email string) (username string, password string, error error) {
	resetCmd := exec.Command(reset_password, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetCmd.Stdout = stdout
	resetCmd.Stderr = stderr

	err := resetCmd.Start()
	if err != nil {
		return "", "", err
	}
	err = resetCmd.Wait()
	if err != nil {
		return "", "", err
	}

	output := strings.Fields(stdout.String())
	username = output[0]
	password = output[1]

	return username, password, nil
}

func CheckIfExists(email string) (result bool, error error) {
	checkIfExistsCmd := exec.Command(check_if_exists, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	checkIfExistsCmd.Stdout = stdout
	checkIfExistsCmd.Stderr = stderr

	err := checkIfExistsCmd.Run()
	if err != nil {
		return false, err
	}

	output := strings.TrimSpace(stdout.String())
	if output == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func SetOpenstackRC() {
	os.Setenv("OS_AUTH_URL", OS_AUTH_URL)
	os.Setenv("OS_PROJECT_ID", OS_PROJECT_ID)
	os.Setenv("OS_PROJECT_NAME", OS_PROJECT_NAME)
	os.Setenv("OS_USER_DOMAIN_NAME", OS_USER_DOMAIN_NAME)
	os.Setenv("OS_PROJECT_DOMAIN_ID", OS_PROJECT_DOMAIN_ID)
	os.Setenv("OS_USERNAME", OS_USERNAME)
	os.Setenv("OS_PASSWORD", OS_PASSWORD)
	os.Setenv("OS_REGION_NAME", OS_REGION_NAME)
	os.Setenv("OS_INTERFACE", OS_INTERFACE)
	os.Setenv("OS_IDENTITY_API_VERSION", OS_IDENTITY_API_VERSION)
}
