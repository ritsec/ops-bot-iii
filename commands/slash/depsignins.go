package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// TODO: update config and get config
	deprecateSigninsChannel string
)

func Depsignins() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "depsignins",
			Description:              "Deprecate the Signins",
			DefaultMemberPermissions: &permission.Admin,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.member:Member",
				tracer.ResourceName("/member"),
			)
			defer span.Finish()

			logging.Debug(s, "Deprecate signins received", i.Member.User, span)

			err := secondUserApproval(s, i, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}
		}
}

func secondUserApproval(s *discordgo.Session, i *discordgo.InteractionCreate, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.slash.depsignins.secondUserApproval",
		tracer.ResourceName("/depsignins:secondUserApproval"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	confirmSlug := uuid.New().String()
	denySlug := uuid.New().String()
	responseChan := make(chan bool)

	(*ComponentHandlers)[confirmSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- true
	}
	defer delete(*ComponentHandlers, confirmSlug)

	(*ComponentHandlers)[denySlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreateChan <- i
		responseChan <- false
	}
	defer delete(*ComponentHandlers, denySlug)

	m, err := s.ChannelMessageSendComplex(deprecateSigninsChannel, &discordgo.MessageSend{
		Content: "Approve deprecating the signins?",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: confirmSlug,
						Label:    "Confirm",
						Style:    discordgo.SuccessButton,
						Emoji: &discordgo.ComponentEmoji{
							Name: "✔️",
						},
					},
					discordgo.Button{
						CustomID: denySlug,
						Label:    "Deny",
						Style:    discordgo.DangerButton,
						Emoji: &discordgo.ComponentEmoji{
							Name: "✖️",
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
		// TODO: Deprecate Signins here
		return nil
	} else {
		_, err = s.ChannelMessageEdit(deprecateSigninsChannel, m.ID, "Deprecate signins denied")
		if err != nil {
			return err
		}
		return nil
	}
}
