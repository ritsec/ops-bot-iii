package scheduled

import (
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// HeartbeatURL is the URL to send the heartbeat to
	heartbeatURL string = config.GetString("commands.heartbeat.url")
)

// Heartbeat is a scheduled task that sends a heartbeat to the push-based life check
func Heartbeat() *structs.ScheduledEvent {
	return structs.NewScheduledTask(
		func(s *discordgo.Session, quit chan interface{}) error {
			ticker := time.NewTicker(10 * time.Second)
			for {
				select {
				case <-quit:
					return nil
				case <-ticker.C:
					span := tracer.StartSpan(
						"commands.scheduled.heartbeat:Heartbeat",
						tracer.ResourceName("Scheduled.Heartbeat"),
					)

					_, err := http.Get(heartbeatURL)
					if err != nil {
						logging.Error(s, err.Error(), nil, span)
					} else {
						logging.DebugLow(s, "Heartbeat sent", nil, span)
					}

					span.Finish()
				}
			}
		},
	)

}
