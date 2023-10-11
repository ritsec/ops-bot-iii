package slash

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/mail"
	"github.com/ritsec/ops-bot-iii/structs"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// Channel ID of the member approval channel
	memberApprovalChannel string = config.GetString("commands.member.channel_id")

	// Role ID of the member role
	memberRole string = config.GetString("commands.member.member_role_id")

	// Role ID of the external role
	externalRole string = config.GetString("commands.member.external_role_id")

	// Role ID of the prospective role
	prosectiveRole string = config.GetString("commands.member.prospective_role_id")

	// Role ID of the staff role
	staffRole string = config.GetString("commands.member.staff_role_id")

	// Role ID of the alumni role
	alumniRole string = config.GetString("commands.member.alumni_role_id")
)

// Member is the member command
func Member() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "member",
			Description: "Become a member through our verification process",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.member:Member",
				tracer.ResourceName("/member"),
			)
			defer span.Finish()

			logging.Debug(s, "Member command received", i.Member.User, span)

			// check if user is already a member
			if data.User.IsVerified(i.Member.User.ID, span.Context()) {
				err := addMemberRole(s, i, "", 0, true, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}
				logging.Debug(s, "User is already verified", i.Member.User, span)
				return
			}

			logging.Debug(s, "User is not verified", i.Member.User, span)

			// check if user already has too many verification attempts
			attempts, err := data.User.GetVerificationAttempts(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			if attempts >= 5 {
				logging.Debug(s, fmt.Sprintf("User has too many verification attempts: %d", attempts), i.Member.User, span)
				err := tooManyAttempts(s, i, attempts, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					return
				}
				return
			}

			logging.Debug(s, fmt.Sprintf("User has %d previous verification attempts", attempts), i.Member.User, span)

			// increment verification attempts
			_, err = data.User.IncrementVerificationAttempts(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			originalInteraction := i

			// check if user has an RIT userEmail
			ritEmail, i, err := hasRITEmail(s, i, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
				return
			}

			if ritEmail {
				// get userEmail
				userEmail, i, err := getEmail(s, i, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
					return
				}

				logging.Debug(s, fmt.Sprintf("User provided email: `%v`", userEmail), originalInteraction.Member.User, span)

				// check if userEmail is valid
				if !validRITEmail(userEmail, span.Context()) {
					logging.Debug(s, fmt.Sprintf("User has invalid RIT email: `%v`", userEmail), originalInteraction.Member.User, span)
					err := invalidRITEmail(s, i, userEmail, attempts, span.Context())
					if err != nil {
						logging.Error(s, "", originalInteraction.Member.User, span, logrus.Fields{"error": err})
					}
					return
				}

				// check if email is already in use
				if data.User.EmailExists(originalInteraction.Member.User.ID, userEmail, span.Context()) {
					logging.Debug(s, fmt.Sprintf("User has already used email: `%v`", userEmail), originalInteraction.Member.User, span)
					err := emailInUse(s, i, userEmail, attempts, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
					}
					return
				}

				// send userEmail
				code, err := mail.SendVerificationEmail(userEmail, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
					return
				}

				logging.Debug(s, fmt.Sprintf("User send Email with verification code: `%v`", code), originalInteraction.Member.User, span)

				// check if userEmail was recieved
				recieved, i, err := recievedEmail(s, i, userEmail, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
					return
				}

				if recieved {

					logging.Debug(s, fmt.Sprintf("User recieved email: `%v`", userEmail), originalInteraction.Member.User, span)

					// get verification code
					verificationCode, i, err := getVerificationCode(s, i, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
						return
					}

					// check code
					if strings.TrimSpace(code) != strings.TrimSpace(verificationCode) {
						logging.Debug(s, "User provided invalid verification code", originalInteraction.Member.User, span)
						err := invalidCode(s, i, verificationCode, attempts, span.Context())
						if err != nil {
							logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
						}
						return
					}

					// add member role
					err = addMemberRole(s, i, userEmail, attempts, false, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), originalInteraction.Member.User, span, logrus.Fields{"error": err})
						return
					}

				} else {
					manualVerification(s, i, userEmail, attempts, span.Context())
					return
				}

			} else {
				manualVerification(s, i, "", attempts, span.Context())
				return
			}
		},
	}
}

// tooManyAttempts is called when a user has too many verification attempts
// it sends a message to the user asking if they want to use to manual verification
func tooManyAttempts(s *discordgo.Session, i *discordgo.InteractionCreate, attempts int, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.slash.member.tooManyAttempts",
		tracer.ResourceName("/member:tooManyAttempts"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	yesSlug := uuid.New().String()
	noSlug := uuid.New().String()

	responseChan := make(chan bool)

	(*ComponentHandlers)[yesSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- true
	}
	defer delete(*ComponentHandlers, yesSlug)

	(*ComponentHandlers)[noSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- false
	}
	defer delete(*ComponentHandlers, noSlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "You have too many verification attemps! Would you like to proceed with manual verification instead?",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Yes",
							Style:    discordgo.SuccessButton,
							CustomID: yesSlug,
						},
						discordgo.Button{
							Label:    "No",
							Style:    discordgo.DangerButton,
							CustomID: noSlug,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	i = <-interactionCreateChan
	response := <-responseChan

	if response {
		manualVerification(s, i, "", attempts, span.Context())
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Verification cancelled.",
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// emailInUse is called when another user has already used the email they provided
// it sends a message to the user asking if they want to use to manual verification
func emailInUse(s *discordgo.Session, i *discordgo.InteractionCreate, userEmail string, attempts int, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.slash.member:emailInUse",
		tracer.ResourceName("/member:emailInUse"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	yesSlug := uuid.New().String()
	noSlug := uuid.New().String()

	responseChan := make(chan bool)

	(*ComponentHandlers)[yesSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- true
	}
	defer delete(*ComponentHandlers, yesSlug)

	(*ComponentHandlers)[noSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- false
	}
	defer delete(*ComponentHandlers, noSlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("The email %s is already in use. Would you like to proceed with manual verification instead.", userEmail),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Yes",
							Style:    discordgo.SuccessButton,
							CustomID: yesSlug,
						},
						discordgo.Button{
							Label:    "No",
							Style:    discordgo.DangerButton,
							CustomID: noSlug,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	i = <-interactionCreateChan
	response := <-responseChan

	if response {
		manualVerification(s, i, userEmail, attempts, span.Context())
		return nil
	} else {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Verification cancelled. Better luck next time!",
			},
		})
	}
}

// addMemberRole adds the member role to the user and marks them as verified
func addMemberRole(s *discordgo.Session, i *discordgo.InteractionCreate, userEmail string, attempts int, firstMessage bool, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.slash.member:addMemberRole",
		tracer.ResourceName("/member:addMemberRole"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, memberRole)
	if err != nil {
		return err
	}

	// add email to user
	_, err = data.User.SetEmail(i.Member.User.ID, userEmail, span.Context())
	if err != nil {
		return err
	}

	// mark user as verified
	_, err = data.User.MarkVerified(i.Member.User.ID, span.Context())
	if err != nil {
		return err
	}

	logging.Debug(s, fmt.Sprintf("User successfully verified:\n Email:`%v`\nAttempts:`%d`", userEmail, attempts), i.Member.User, span)
	if firstMessage {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "You have been verified as a member of RITSEC. Welcome!",
			},
		})
	} else {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "You have been verified as a member of RITSEC. Welcome!",
			},
		})
	}
}

// getVerificationCode sends the user a verification code and waits for them to respond with it
// if the user fails to respond with the code, they are given the option to resend the code or cancel verification
func getVerificationCode(s *discordgo.Session, i *discordgo.InteractionCreate, ctx ddtrace.SpanContext) (string, *discordgo.InteractionCreate, error) {
	span := tracer.StartSpan(
		"commands.slash.member:getVerificationCode",
		tracer.ResourceName("/member:getVerificationCode"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	verifyChan := make(chan string)
	defer close(verifyChan)

	verifySlug := uuid.New().String()

	(*ComponentHandlers)[verifySlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ModalSubmitData()

		verifyChan <- data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, verifySlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: verifySlug,
			Title:    "Member verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "verify",
							Label:       "What is the verification code?",
							Style:       discordgo.TextInputShort,
							Placeholder: "000000",
							Required:    true,
							MaxLength:   6,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return "", nil, err
	}

	verificationCode := <-verifyChan
	i = <-interactionCreateChan

	return verificationCode, i, nil
}

// recievedEmail will prompt a user to check if they recieved an email
func recievedEmail(s *discordgo.Session, i *discordgo.InteractionCreate, userEmail string, ctx ddtrace.SpanContext) (bool, *discordgo.InteractionCreate, error) {
	span := tracer.StartSpan(
		"commands.slash.member:recievedEmail",
		tracer.ResourceName("/member:recievedEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	recievedSlug := uuid.New().String()
	unrecievedSlug := uuid.New().String()

	recievedChan := make(chan bool)
	defer close(recievedChan)

	(*ComponentHandlers)[recievedSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		recievedChan <- true
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, recievedSlug)

	(*ComponentHandlers)[unrecievedSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		recievedChan <- false
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, unrecievedSlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "A verification code has been sent to \"" + userEmail + "\". This could take up to 10 minutes. **Do not close discord or this window will be closed**.",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: recievedSlug,
							Label:    "I recieved the code",
							Style:    discordgo.SuccessButton,
						},
						discordgo.Button{
							CustomID: unrecievedSlug,
							Label:    "I did not recieve the code",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return false, nil, err
	}

	recieved := <-recievedChan
	i = <-interactionCreateChan

	return recieved, i, nil
}

// invalidCode is called when a user enters an invalid code
// they are given the option to try again, cancel verification, or preform manual verification
func invalidCode(s *discordgo.Session, i *discordgo.InteractionCreate, userEmail string, attempts int, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.slash.member:invalidCode",
		tracer.ResourceName("/member:invalidCode"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	continueSlug := uuid.New().String()
	quitSlug := uuid.New().String()

	responseChan := make(chan bool)

	(*ComponentHandlers)[continueSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- true
	}
	defer delete(*ComponentHandlers, continueSlug)

	(*ComponentHandlers)[quitSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- false
	}
	defer delete(*ComponentHandlers, quitSlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "The code you returned is not valid. Would you like to continue with manual verification?",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: continueSlug,
							Label:    "yes",
							Style:    discordgo.SuccessButton,
						},
						discordgo.Button{
							CustomID: quitSlug,
							Label:    "no",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	i = <-interactionCreateChan
	response := <-responseChan

	if !response {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Verification attempt cancelled. Better luck next time!",
			},
		})
	}

	manualVerification(s, i, userEmail, attempts, span.Context())
	return nil
}

// invalidRITEmail is called when a user enters an invalid RIT email
// they are given the option to try again, cancel verification, or preform manual verification
func invalidRITEmail(s *discordgo.Session, i *discordgo.InteractionCreate, userEmail string, attempts int, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.slash.member:invalidRITEmail",
		tracer.ResourceName("/member:invalidRITEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	continueSlug := uuid.New().String()
	quitSlug := uuid.New().String()

	responseChan := make(chan bool)

	(*ComponentHandlers)[continueSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- true
	}
	defer delete(*ComponentHandlers, continueSlug)

	(*ComponentHandlers)[quitSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- false
	}
	defer delete(*ComponentHandlers, quitSlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("The email address `%s` is not a valid RIT email address. Would you like to continue with manual verification?", userEmail),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: continueSlug,
							Label:    "yes",
							Style:    discordgo.SuccessButton,
						},
						discordgo.Button{
							CustomID: quitSlug,
							Label:    "no",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	i = <-interactionCreateChan
	response := <-responseChan

	if !response {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Verification attempt cancelled. Better luck next time!",
			},
		})
	}

	manualVerification(s, i, userEmail, attempts, span.Context())
	return nil
}

// validRITEmail is used to check if an email is a valid RIT email
func validRITEmail(userEmail string, ctx ddtrace.SpanContext) bool {
	span := tracer.StartSpan(
		"commands.slash.member:validRITEmail",
		tracer.ResourceName("/member:validRITEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	ritEmailRegex := `^[a-z]{2,3}\d{4}@(g\.)?rit\.edu$`
	matched, err := regexp.MatchString(ritEmailRegex, userEmail)
	if err != nil {
		return false
	}
	return matched
}

// getEmail is used to get the email of a user
func getEmail(s *discordgo.Session, i *discordgo.InteractionCreate, ctx ddtrace.SpanContext) (string, *discordgo.InteractionCreate, error) {
	span := tracer.StartSpan(
		"commands.slash.member:getEmail",
		tracer.ResourceName("/member:getEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	email_slug := uuid.New().String()

	emailChan := make(chan string)
	interactionCreateChan := make(chan *discordgo.InteractionCreate)

	defer close(emailChan)
	defer close(interactionCreateChan)

	(*ComponentHandlers)[email_slug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ModalSubmitData()

		emailChan <- data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		interactionCreateChan <- i
	}

	defer delete(*ComponentHandlers, email_slug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: email_slug,
			Title:    "Member verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "email",
							Label:       "What is your @rit.edu email address?",
							Style:       discordgo.TextInputShort,
							Placeholder: "abc1234@rit.edu",
							Required:    true,
							MaxLength:   18,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return "", nil, err
	}

	userEmail := <-emailChan
	i = <-interactionCreateChan

	return userEmail, i, nil
}

// hasRITEmail is used to check if a user has a valid RIT email
func hasRITEmail(s *discordgo.Session, i *discordgo.InteractionCreate, ctx ddtrace.SpanContext) (bool, *discordgo.InteractionCreate, error) {
	span := tracer.StartSpan(
		"commands.slash.member:hasRITEmail",
		tracer.ResourceName("/member:hasRITEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	emailChan := make(chan bool)
	defer close(emailChan)

	yesSlug := uuid.New().String()
	noSlug := uuid.New().String()

	(*ComponentHandlers)[yesSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		emailChan <- true
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, yesSlug)

	(*ComponentHandlers)[noSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		emailChan <- false
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, noSlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Do you have an RIT email address?",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: yesSlug,
							Label:    "Yes",
							Style:    discordgo.SuccessButton,
						},
						discordgo.Button{
							CustomID: noSlug,
							Label:    "No",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	})

	if err != nil {
		return false, nil, err
	}

	userEmail := <-emailChan
	i = <-interactionCreateChan

	return userEmail, i, nil
}

// manualVerification is used to manually verify a user
// This is used when the user fail other verification methods
func manualVerification(s *discordgo.Session, i *discordgo.InteractionCreate, userEmail string, attempts int, ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.slash.member:manualVerification",
		tracer.ResourceName("/member:manualVerification"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	user := i.Member.User

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	verifyChan := make(chan string)
	defer close(verifyChan)

	verifySlug := uuid.New().String()

	(*ComponentHandlers)[verifySlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ModalSubmitData()

		verifyChan <- data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, verifySlug)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: verifySlug,
			Title:    "Member verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "verify",
							Label:    "Who are you? Why would you like to join?",
							Style:    discordgo.TextInputParagraph,
							Required: true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
		return
	}

	message := <-verifyChan
	i = <-interactionCreateChan

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Message has been sent to RITSEC E-Board for verification. Please wait for a response.",
		},
	})
	if err != nil {
		logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
		return
	}

	memberSlug := uuid.New().String()
	externalSlug := uuid.New().String()
	prosectiveSlug := uuid.New().String()
	staffSlug := uuid.New().String()
	alumniSlug := uuid.New().String()
	denySlug := uuid.New().String()

	memberChan := make(chan string)
	defer close(memberChan)

	(*ComponentHandlers)[memberSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		memberChan <- "member"
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, memberSlug)

	(*ComponentHandlers)[externalSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		memberChan <- "external"
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, externalSlug)

	(*ComponentHandlers)[prosectiveSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		memberChan <- "prosective"
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, prosectiveSlug)

	(*ComponentHandlers)[staffSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		memberChan <- "staff"
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, staffSlug)

	(*ComponentHandlers)[alumniSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		memberChan <- "alumni"
		interactionCreateChan <- i
	}

	(*ComponentHandlers)[denySlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		memberChan <- "deny"
		interactionCreateChan <- i
	}
	defer delete(*ComponentHandlers, denySlug)

	m, err := s.ChannelMessageSendComplex(memberApprovalChannel, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Verification request",
			Fields: []*discordgo.MessageEmbedField{
				func() *discordgo.MessageEmbedField {
					if userEmail == "" {
						return &discordgo.MessageEmbedField{
							Name:  "email",
							Value: "Not provided",
						}
					} else {
						return &discordgo.MessageEmbedField{
							Name:  "email",
							Value: userEmail,
						}
					}
				}(),
				{
					Name:  "Discord",
					Value: user.Mention(),
				},
				{
					Name:  "Message",
					Value: message,
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: memberSlug,
						Label:    "Member",
						Style:    discordgo.SuccessButton,
					},
					discordgo.Button{
						CustomID: externalSlug,
						Label:    "External",
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: prosectiveSlug,
						Label:    "Prospective",
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: staffSlug,
						Label:    "Staff",
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: alumniSlug,
						Label:    "Alumni",
						Style:    discordgo.PrimaryButton,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: denySlug,
						Label:    "Deny",
						Style:    discordgo.DangerButton,
					},
				},
			},
		},
	})
	if err != nil {
		logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
		return
	}

	memberType := <-memberChan
	i = <-interactionCreateChan

	switch memberType {
	case "member":
		err = addMemberRole(s, i, userEmail, attempts, false, span.Context())
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}
	case "external":
		err = s.GuildMemberRoleAdd(i.GuildID, user.ID, externalRole)
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}

		err = helpers.SendDirectMessage(s, user.ID, "You have been verified as an external member of RITSEC. Welcome!", span.Context())
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}
	case "prosective":
		err = s.GuildMemberRoleAdd(i.GuildID, user.ID, prosectiveRole)
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}

		err = helpers.SendDirectMessage(s, user.ID, "You have been verified as a prospective member of RITSEC. Welcome!", span.Context())
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}
	case "staff":
		err = s.GuildMemberRoleAdd(i.GuildID, user.ID, staffRole)
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}

		err = helpers.SendDirectMessage(s, user.ID, "You have been verified as a staff member of RIT. Welcome!", span.Context())
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}
	case "alumni":
		err = s.GuildMemberRoleAdd(i.GuildID, user.ID, alumniRole)
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}

		err = helpers.SendDirectMessage(s, user.ID, "You have been verified as an alumni of RITSEC. Welcome!", span.Context())
		if err != nil {
			logging.Error(s, err.Error(), user, span, logrus.Fields{"error": err})
			return
		}
	}

	err = s.ChannelMessageDelete(memberApprovalChannel, m.ID)
	if err != nil {
		logging.Error(s, "Error encounted while deleting channel message", user, span, logrus.Fields{"error": err})
		return
	}
}
