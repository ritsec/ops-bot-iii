package slash

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Import slash command
func Import() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "import",
			Description:              "Import signins from csv",
			DefaultMemberPermissions: &permission.Admin,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "file",
					Description: "A csv file containing sign-ins to import.",
					Required:    true,
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.import:Import",
				tracer.ResourceName("/import"),
			)
			defer span.Finish()

			logging.Debug(s, "Import command received", i.Member.User, span)

			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Reading CSV file...",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				},
			)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
				return
			}

			data := i.ApplicationCommandData()
			var attachment *discordgo.MessageAttachment
			if option := data.GetOption("file"); option != nil && option.Type == discordgo.ApplicationCommandOptionAttachment && data.Resolved != nil {
				if attachmentID, ok := option.Value.(string); ok {
					attachment = data.Resolved.Attachments[attachmentID]
				}
			}

			rows, err := Parse_csv(attachment)
			var message string
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
				message = "Could not parse CSV: " + err.Error()
			} else {
				message = fmt.Sprintf("Parsed %d CSV rows (including the header, if present).", len(rows))
			}

			var message2 strings.Builder
			for _, value := range rows[0] {
				message2.WriteString(value + " ")
			}
			logging.Debug(s, message2.String(), i.Member.User, span)

			if err := helpers.IntRespondEdit(s, i, message); err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}
		}
}

// Parse_csv downloads an attachment and reads its CSV records, preserving
// the first row so the caller can decide whether it contains column headers.
func Parse_csv(file *discordgo.MessageAttachment) ([][]string, error) {
	if file == nil || file.URL == "" {
		return nil, fmt.Errorf("missing CSV attachment")
	}

	const maxSize = 10 * 1024 * 1024
	if file.Size > maxSize {
		return nil, fmt.Errorf("CSV file exceeds the 10 MiB limit")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Get(file.URL)
	if err != nil {
		return nil, fmt.Errorf("download CSV: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download CSV: HTTP %d", response.StatusCode)
	}

	// Bound the download even if the attachment's reported size is incorrect.
	limited := &io.LimitedReader{R: response.Body, N: maxSize + 1}
	rows, err := csv.NewReader(limited).ReadAll()
	if limited.N == 0 {
		return nil, fmt.Errorf("CSV file exceeds the 10 MiB limit")
	}
	if err != nil {
		return nil, fmt.Errorf("parse CSV: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	return rows, nil
}
