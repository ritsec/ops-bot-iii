package slash

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/mail"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Verify() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "verify",
			Description: "Verify your email to use services like openstack and count for attendance",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.verify:Verify",
				tracer.ResourceName("/Verify"),
			)
			defer span.Finish()

			logging.Debug(s, "Verify command received", i.Member.User, span)
			helpers.InitialMessage(s, i, "Checking to see if you are verified...")

			if data.User.IsVerified(i.Member.User.ID, span.Context()) {
				logging.Debug(s, "User is already verified", i.Member.User, span)
				helpers.UpdateMessage(s, i, "You are already verified.")
				return
			}

			// check if user has an RIT userEmail
			ritEmail, i, err := hasRITEmail(s, i, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			if ritEmail {
				// get userEmail
				userEmail, i, err := getEmail(s, i, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}

				logging.Debug(s, fmt.Sprintf("User provided email: `%v`", userEmail), i.Member.User, span)

				// check if userEmail is valid
				if !validRITEmail(userEmail, span.Context()) {
					logging.Debug(s, fmt.Sprintf("User has invalid RIT email: `%v`", userEmail), i.Member.User, span)
					helpers.UpdateMessage(s, i, "Invalid RIT email.")
					return
				}

				// check if email is already in use
				if data.User.EmailExists(i.Member.User.ID, userEmail, span.Context()) {
					logging.Debug(s, fmt.Sprintf("User has already used email: `%v`", userEmail), i.Member.User, span)
					helpers.UpdateMessage(s, i, "Email already in use.")
					return
				}

				// send userEmail
				helpers.UpdateMessage(s, i, "Sending a verification email...")
				code, err := mail.SendVerificationEmail(userEmail, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}

				logging.Debug(s, fmt.Sprintf("User send Email with verification code: `%v`", code), i.Member.User, span)

				// check if userEmail was recieved
				recieved, i, err := recievedEmail(s, i, userEmail, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}

				if recieved {

					logging.Debug(s, fmt.Sprintf("User recieved email: `%v`", userEmail), i.Member.User, span)

					// get verification code
					verificationCode, i, err := getVerificationCode(s, i, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
						return
					}

					// check code
					if strings.TrimSpace(code) != strings.TrimSpace(verificationCode) {
						logging.Debug(s, "User provided invalid verification code", i.Member.User, span)
						err := invalidCode(s, i, verificationCode, 0, span.Context())
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
						}
						return
					}

					// add email to user
					_, err = data.User.SetEmail(i.Member.User.ID, userEmail, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
						return
					}

					// mark user as verified
					_, err = data.User.MarkVerified(i.Member.User.ID, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
						return
					}

					msg := "You have been verified as a member of RITSEC. Welcome!"
					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Content: &msg,
					})
					if err != nil {
						logging.Error(s, err.Error(), i.User, span, logrus.Fields{"error": err})
						return
					}
				} else {
					helpers.UpdateMessage(s, i, "Verification canceled.")
				}
			}
		}
}
