package scheduled

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// HeartbeatProvider is the heartbeat provider, either "kuma" or "datadog"
	heartbeatProvider string = config.GetString("commands.heartbeat.provider")

	// HeartbeatURL is the URL to send the heartbeat to (kuma)
	heartbeatURL string = config.GetString("commands.heartbeat.url")

	// Datadog heartbeat settings
	ddAPIKey   string = config.GetString("commands.heartbeat.datadog.api_key")
	ddAppKey   string = config.GetString("commands.heartbeat.datadog.app_key")
	ddSite     string = config.GetString("commands.heartbeat.datadog.site")
	ddHostname string = config.GetString("commands.heartbeat.datadog.hostname")
)

// Heartbeat is a scheduled task that sends a heartbeat to the configured push-based life check
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

			if err := sendHeartbeat(); err != nil {
				logging.Error(s, err.Error(), nil, span)
			} else {
				logging.DebugLow(s, "Heartbeat sent", nil, span)
			}

			span.Finish()
		}
	}
}

// sendHeartbeat sends a heartbeat to the configured provider
func sendHeartbeat() error {
	switch heartbeatProvider {
	case "datadog":
		return sendDatadogHeartbeat()
	default:
		return sendKumaHeartbeat()
	}
}

// sendKumaHeartbeat performs an HTTP GET against the Uptime Kuma push-based life check URL
func sendKumaHeartbeat() error {
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Get(heartbeatURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Drain body to allow connection reuse.
	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		return copyErr
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("heartbeat returned non-2xx status: %s", resp.Status)
	}
	return nil
}

// sendDatadogHeartbeat POSTs a heartbeat to the Datadog push-based heartbeat API
func sendDatadogHeartbeat() error {
	site := ddSite
	if site == "" {
		site = "datadoghq.com"
	}

	body, err := json.Marshal(map[string]string{
		"hostname": ddHostname,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.%s/api/v1/heartbeat", site),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", ddAPIKey)
	req.Header.Set("DD-APPLICATION-KEY", ddAppKey)

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Drain body to allow connection reuse.
	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		return copyErr
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog heartbeat returned non-2xx status: %s", resp.Status)
	}
	return nil
}
