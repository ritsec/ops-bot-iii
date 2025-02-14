package slash

import (
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
	// FeedbackChannelID is the ID of the channel to send feedback to
	FeedbackChannelID string = config.GetString("commands.feedback.channel_id")
)

// Feedback is the feedback command
func Feedback() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "feedback",
			Description:              "Send Anonymous Feedback to E-Board",
			DefaultMemberPermissions: &permission.Member,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.feedback:Feedback",
				tracer.ResourceName("/feedback"),
			)
			defer span.Finish()

			logging.Debug(s, "Feedback command received", nil, span)

			feedbackSlug := uuid.New().String()

			closeChan := make(chan bool)
			defer close(closeChan)

			(*ComponentHandlers)[feedbackSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				span_feedbackSlug := tracer.StartSpan(
					"commands.slash.feedback:Feedback:feedbackSlug",
					tracer.ResourceName("/feedback:feedbackSlug"),
					tracer.ChildOf(span.Context()),
				)
				defer span_feedbackSlug.Finish()

				feedback := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Feedback Sent",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					// Let the eboard know that this failed while keeping the user anon
					logging.Error(s, err.Error(), nil, span_feedbackSlug, logrus.Fields{"error": err})

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Unable to send the feedback. Error reported without revealing your name.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})

					if err != nil {
						logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
						err = helpers.SendDirectMessage(s, i.Member.User.ID, "Unable to send the feedback. Error reported without revealing your name.", span.Context())

						// dawg everything failed, might as well have eboard contact the user
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
						}
					}
				}

				_, err = s.ChannelMessageSendComplex(FeedbackChannelID,
					&discordgo.MessageSend{
						Embed: &discordgo.MessageEmbed{
							Title:       "New Feedback",
							Description: feedback,
						},
					},
				)
				if err != nil {
					// Let the eboard know that this failed while keeping the user anon
					logging.Error(s, err.Error(), nil, span_feedbackSlug, logrus.Fields{"error": err})

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Unable to send the feedback. Error reported without revealing your name.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})

					if err != nil {
						logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
						err = helpers.SendDirectMessage(s, i.Member.User.ID, "Unable to send the feedback. Error reported without revealing your name.", span.Context())

						// dawg everything failed, might as well have eboard contact the user
						if err != nil {
							logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
						}
					}
				}

				closeChan <- true
			}

			defer delete(*ComponentHandlers, feedbackSlug)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: feedbackSlug,
					Title:    "Feedback",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID: "feedback",
									Label:    "Feedback",
									Style:    discordgo.TextInputParagraph,
								},
							},
						},
					},
				},
			})
			if err != nil {
				// Let the eboard know that this failed while keeping the user anon
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Unable to send the feedback. Error reported without revealing your name.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})

				if err != nil {
					logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
					err = helpers.SendDirectMessage(s, i.Member.User.ID, "Unable to send the feedback. Error reported without revealing your name.", span.Context())

					// dawg everything failed, might as well have eboard contact the user
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
					}
				}
			}
			logging.Debug(s, "Feedback command responded", nil, span)

			<-closeChan

			logging.Debug(s, "Feedback command closed", nil, span)
		}
}
