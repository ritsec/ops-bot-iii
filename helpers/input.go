package helpers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// SelectButtons sends a message with buttons to select from
func SelectButtons(s *discordgo.Session, i *discordgo.InteractionCreate, componentHandlers *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), ctx ddtrace.SpanContext, message string, options ...string) (string, *discordgo.InteractionCreate, error) {
	span := tracer.StartSpan(
		"helpers.input:SelectButtons",
		tracer.ResourceName("helpers.input:SelectButtons"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	components := []discordgo.MessageComponent{}
	slugs := []string{}
	selected := make(chan string)
	interaction := make(chan *discordgo.InteractionCreate)

	defer close(selected)
	defer close(interaction)

	emojis := []string{
		"🟥",
		"🟧",
		"🟨",
		"🟩",
		"🟦",
		"🟪",
		"⬛",
		"⬜",
		"🟥",
		"🟧",
		"🟨",
		"🟩",
		"🟦",
		"🟪",
		"⬛",
		"⬜",
	}

	for j, option := range options {
		slug := uuid.New().String()
		components = append(components, discordgo.Button{
			Label:    option,
			Style:    discordgo.PrimaryButton,
			CustomID: slug,
			Emoji: discordgo.ComponentEmoji{
				Name: emojis[j],
			},
		})
		slugs = append(slugs, slug)

		(*componentHandlers)[slug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			selected <- i.MessageComponentData().CustomID
			interaction <- i
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	})

	if err != nil {
		return "", nil, err
	}

	choice := <-selected
	newInteraction := <-interaction

	for _, slug := range slugs {
		delete(*componentHandlers, slug)
	}
	return options[IndexOf(slugs, choice)], newInteraction, nil
}

// SelectButtonsEdit edits a message with buttons to select from
func SelectButtonsEdit(s *discordgo.Session, i *discordgo.InteractionCreate, componentHandlers *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), ctx ddtrace.SpanContext, message string, options ...string) (string, error) {
	span := tracer.StartSpan(
		"helpers.input:SelectButtonsEdit",
		tracer.ResourceName("helpers.input:SelectButtonsEdit"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	components := []discordgo.MessageComponent{}
	slugs := []string{}
	selected := make(chan string)
	defer close(selected)

	emojis := []string{
		"🟥",
		"🟧",
		"🟨",
		"🟩",
		"🟦",
		"🟪",
		"⬛",
		"⬜",
		"🟥",
		"🟧",
		"🟨",
		"🟩",
		"🟦",
		"🟪",
		"⬛",
		"⬜",
	}

	for j, option := range options {
		slug := uuid.New().String()
		components = append(components, discordgo.Button{
			Label:    option,
			Style:    discordgo.PrimaryButton,
			CustomID: slug,
			Emoji: discordgo.ComponentEmoji{
				Name: emojis[j],
			},
		})
		slugs = append(slugs, slug)

		(*componentHandlers)[slug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			selected <- i.MessageComponentData().CustomID
		}
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &message,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: components,
			},
		},
	})

	if err != nil {
		return "", err
	}

	choice := <-selected

	for _, slug := range slugs {
		delete(*componentHandlers, slug)
	}
	return options[IndexOf(slugs, choice)], nil
}
