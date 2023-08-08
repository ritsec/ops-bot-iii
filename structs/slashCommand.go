package structs

import "github.com/bwmarrin/discordgo"

// SlashCommand is a struct that represents a slash command
type SlashCommand struct {
	// Command is the command to register
	Command *discordgo.ApplicationCommand

	// Handler is the handler for the command
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
