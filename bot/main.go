package bot

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
)

var (
	// Session is the global discordgo session
	Session *discordgo.Session
)

func init() {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	// Set the global session
	Session = session

	// Set the intents
	Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// This is required to use handlers like MessageEdit
	// Set the max message count
	Session.State.MaxMessageCount = 50
}
