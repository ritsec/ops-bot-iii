package helpers

import (
	"github.com/bwmarrin/discordgo"
)

func IntRespondEdit(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &message,
	})
	if err != nil {
		return err
	}
	return nil
}
