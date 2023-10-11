package slash

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Reboot is a slash command that reboots the bot
func Reboot() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:                     "reboot",
			Description:              "Reboot the bot",
			DefaultMemberPermissions: &permission.Admin,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.reboot:Reboot",
				tracer.ResourceName("/reboot"),
			)
			defer span.Finish()

			logging.Debug(s, "Reboot command received", i.Member.User, span)
			logging.Critical(s, "Rebooting bot", i.Member.User, span)

			span.Finish()

			os.Exit(0)
		},
	}
}
