package helpers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
)

// JumpURL returns the URL to jump to a message
func JumpURL(m *discordgo.Message) string {
	return "https://discordapp.com/channels/" + config.GuildID + "/" + m.ChannelID + "/" + m.ID
}
