package scheduled

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// updateStatus updates the bot's status to a random activity
func updateStatus(s *discordgo.Session, ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.scheduled.status:updateStatus",
		tracer.ResourceName("Scheduled.Status:updateStatus"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return
	}

	now := time.Now().In(loc)
	weekday := now.Weekday()
	hour := now.Hour()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if weekday == time.Friday && hour >= 12 && hour < 16 {
		err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "RITSEC General Meeting",
					Type: discordgo.ActivityTypeStreaming,
					URL:  "https://www.twitch.tv/ritsec",
				},
			},
		})
		if err != nil {
			logging.Error(s, err.Error(), nil, span)
		}
	} else {
		err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				&activities[r.Intn(len(activities))],
			},
		})
		if err != nil {
			logging.Error(s, err.Error(), nil, span)
		}
	}
}

// Status is a scheduled task that updates the bot's status
func Status() *structs.ScheduledEvent {
	return structs.NewScheduledTask(
		func(s *discordgo.Session, quit chan interface{}) error {
			span := tracer.StartSpan(
				"commands.scheduled.status:Status",
				tracer.ResourceName("Scheduled.Status"),
			)
			defer span.Finish()

			c := cron.New()

			err := c.AddFunc("0 0 * * * *", func() { updateStatus(s, span.Context()) })
			if err != nil {
				return err
			}

			updateStatus(s, span.Context())

			c.Start()
			<-quit
			c.Stop()

			return nil
		},
	)

}

// activities is a list of activities that the bot can be doing
var activities = []discordgo.Activity{
	{
		Name: "with my feelings",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "nvim",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "vim",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "nano",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "emacs",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "with the kernel",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "with its own code",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "with the source code",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Darknet Diaries",
		Type: discordgo.ActivityTypeListening,
	},
	{
		Name: "The CyberWire",
		Type: discordgo.ActivityTypeListening,
	},
	{
		Name: "the worst shitposts",
		Type: discordgo.ActivityTypeListening,
	},
	{
		Name: "ippsec",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "John Hammond",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "the Matrix",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "your command history",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "Synthwave",
		Type: discordgo.ActivityTypeListening,
	},
	{
		Name: "Hacker Music",
		Type: discordgo.ActivityTypeListening,
	},
	{
		Name: "Napolean Dynamite",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "National Lampoon's Christmas Vacation",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "Mr. Robot",
		Type: discordgo.ActivityTypeWatching,
	},
	{
		Name: "with the Wire",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: ":skull: Applying to Co-Ops",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: ":skull: Applying to Internships",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: ":skull: Applying to Full-Time",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: "TryHackMe",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "HackTheBox",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "OverTheWire",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "with Lockpicks",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Competeting in a CTF",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: "SSH",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Tor",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Wireshark",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Nmap",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Metasploit",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "John the Ripper",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Hashcat",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Burp Suite",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Kubernetes",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Docker",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "with the Cloud",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Windows AD",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "WireGuard",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Installing Arch",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: "Installing Kali",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: "Installing Ubuntu",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: "Installing Fedora",
		Type: discordgo.ActivityTypeCustom,
	},
	{
		Name: "The Stack is Down",
		Type: discordgo.ActivityTypeCustom,
	},
}
