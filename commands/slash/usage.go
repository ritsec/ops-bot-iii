package slash

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/slash/permission"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

// Usage is a slash command that returns system metrics
func Usage() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:                     "usage",
			Description:              "System metrics and usage",
			DefaultMemberPermissions: &permission.Admin,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.usage",
				tracer.ResourceName("/usage"),
			)
			defer span.Finish()

			logging.Debug(s, "Usage command received", i.Member.User, span)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Calculating...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			memory, err := memory.Get()
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
			}

			before, err := cpu.Get()
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
			}
			time.Sleep(time.Duration(5) * time.Second)
			after, err := cpu.Get()
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
			}
			total := float64(after.Total - before.Total)

			empty := ""

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &empty,
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title: "System Metrics",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "CPU Usage",
								Value: fmt.Sprintf("%f %%", float64(after.System-before.System)/total*100),
							},
							{
								Name:  "Memory Usage",
								Value: fmt.Sprintf("%f %%", float64(memory.Used)/float64(memory.Total)*100),
							},
							{
								Name:  "Memory Used",
								Value: fmt.Sprintf("%f GiB", float64(memory.Used)/1024/1024/1024),
							},
							{
								Name:  "Memory Free",
								Value: fmt.Sprintf("%f GiB", float64(memory.Free)/1024/1024/1024),
							},
							{
								Name:  "Memory Total",
								Value: fmt.Sprintf("%f GiB", float64(memory.Total)/1024/1024/1024),
							},
						},
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}
		},
	}
}
