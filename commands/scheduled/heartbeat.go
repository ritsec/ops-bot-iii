package scheduled

import (
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// HeartbeatURL is the URL to send the heartbeat to
	heartbeatURL string = config.GetString("commands.heartbeat.url")
)

// Heartbeat is a scheduled task that sends a heartbeat to the push-based life check
func Heartbeat(s *discordgo.Session, quit chan interface{}) error {
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

			resp, err := http.Get(heartbeatURL)
			if err != nil {
				logging.Error(s, err.Error(), nil, span)
			} else {
				// Drain body to allow connection reuse.
				_, copyErr := io.Copy(io.Discard, resp.Body)
				closeErr := resp.Body.Close()

				if copyErr != nil {
					logging.Error(s, copyErr.Error(), nil, span)
				} else if closeErr != nil {
					logging.Error(s, closeErr.Error(), nil, span)
				} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					logging.Error(s, "Heartbeat returned non-2xx status: "+resp.Status, nil, span)
				} else {
					logging.DebugLow(s, "Heartbeat sent", nil, span)
				}
			}

			span.Finish()
		}
	}
}
