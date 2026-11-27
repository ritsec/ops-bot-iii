package slash

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// importEmailRegex validates a single email address (local part capped at 64
// chars by the regex; total length is additionally checked against 254).
var importEmailRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._%+-]{0,63}@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)

const (
	importMaxRows      = 1000
	importMaxListItems = 10

	// csvManualDeleteNote is appended when the uploaded CSV message could not
	// be deleted automatically.
	csvManualDeleteNote = "\nNote: I could not delete the CSV message; please delete it manually."
)

// Import slash command
func Import() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "import",
			Description:              "Import signins from a CSV of emails for a specific meeting",
			DefaultMemberPermissions: &permission.Admin,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "file",
					Description: "CSV file with one email address per line.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "The meeting type these signins are for.",
					Required:    true,
					Choices:     helpers.SigninTypeChoices(),
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date",
					Description: "The date of the meeting (YYYY-MM-DD or M/D/YYYY).",
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
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			cmd := i.ApplicationCommandData()

			var attachment *discordgo.MessageAttachment
			if option := cmd.GetOption("file"); option != nil && option.Type == discordgo.ApplicationCommandOptionAttachment && cmd.Resolved != nil {
				if attachmentID, ok := option.Value.(string); ok {
					attachment = cmd.Resolved.Attachments[attachmentID]
				}
			}

			stringOption := func(name string) string {
				if option := cmd.GetOption(name); option != nil && option.Type == discordgo.ApplicationCommandOptionString {
					if value, ok := option.Value.(string); ok {
						return value
					}
				}
				return ""
			}
			signinTypeString := stringOption("type")
			dateString := stringOption("date")

			var meetingDate time.Time
			for _, layout := range []string{"2006-01-02", "1/2/2006"} {
				if parsed, parseErr := time.Parse(layout, dateString); parseErr == nil {
					meetingDate = parsed
					break
				}
			}
			if meetingDate.IsZero() {
				logging.Error(s, "Could not parse date: "+dateString, i.Member.User, span, logrus.Fields{"date": dateString})
				if editErr := helpers.IntRespondEdit(s, i, fmt.Sprintf("Could not parse date `%s`. Use YYYY-MM-DD or M/D/YYYY.", dateString)); editErr != nil {
					logging.Error(s, editErr.Error(), i.Member.User, span, logrus.Fields{"error": editErr})
				}
				return
			}

			rows, err := parseCsv(attachment)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				if editErr := helpers.IntRespondEdit(s, i, "Could not parse CSV: "+err.Error()); editErr != nil {
					logging.Error(s, editErr.Error(), i.Member.User, span, logrus.Fields{"error": editErr})
				}
				return
			}

			// Delete the uploaded CSV message so it does not remain in the channel.
			csvDeleted := true
			if deleteErr := deleteCSVMessage(s, i, attachment); deleteErr != nil {
				csvDeleted = false
				logging.Error(s, "Failed to delete CSV message", i.Member.User, span, logrus.Fields{"error": deleteErr})
			}

			if len(rows) > importMaxRows {
				logging.Error(s, fmt.Sprintf("CSV contains %d rows; maximum is %d", len(rows), importMaxRows), i.Member.User, span,
					logrus.Fields{"rows": len(rows)})
				if editErr := helpers.IntRespondEdit(s, i, fmt.Sprintf("CSV contains %d rows; the maximum is %d.", len(rows), importMaxRows)); editErr != nil {
					logging.Error(s, editErr.Error(), i.Member.User, span, logrus.Fields{"error": editErr})
				}
				return
			}

			// Excel exports include a UTF-8 BOM on the first field of the first row.
			rows[0][0] = strings.TrimPrefix(rows[0][0], "\ufeff")

			var (
				validEmails  []string
				invalidRows  []string // "line N: snippet" details
				invalidLines []int    // 1-based row numbers of invalid rows
				seen         = map[string]bool{}
			)
			for rowIdx, row := range rows {
				line := rowIdx + 1
				if len(row) != 1 {
					invalidRows = append(invalidRows, fmt.Sprintf("line %d: expected exactly one field, got %d", line, len(row)))
					invalidLines = append(invalidLines, line)
					continue
				}
				email := strings.TrimSpace(row[0])
				if email == "" {
					continue
				}
				if !importEmailRegex.MatchString(email) || len(email) > 254 {
					snippet := email
					if len(snippet) > 32 {
						snippet = snippet[:32] + "…"
					}
					invalidRows = append(invalidRows, fmt.Sprintf("line %d: %q", line, snippet))
					invalidLines = append(invalidLines, line)
					continue
				}
				key := strings.ToLower(email)
				if seen[key] {
					continue
				}
				seen[key] = true
				validEmails = append(validEmails, email)
			}

			if len(validEmails) == 0 {
				logging.Error(s, "No valid emails found in CSV", i.Member.User, span, logrus.Fields{"invalidRows": len(invalidLines)})
				message := "No valid email addresses were found in the CSV."
				if len(invalidRows) > 0 {
					message += fmt.Sprintf("\nInvalid rows (%d): %s", len(invalidRows), cappedList(invalidRows, importMaxListItems))
				}
				if !csvDeleted {
					message += csvManualDeleteNote
				}
				if editErr := helpers.IntRespondEdit(s, i, message); editErr != nil {
					logging.Error(s, editErr.Error(), i.Member.User, span, logrus.Fields{"error": editErr})
				}
				return
			}

			entType := helpers.StringToType(signinTypeString)

			var (
				importedCount int
				alreadyList   []string
				notFoundList  []string
				lookupErrors  int
			)
			for _, email := range validEmails {
				user, err := data.User.FindByEmail(email, span.Context())
				if err != nil {
					logging.Error(s, "Failed to look up user: "+email, i.Member.User, span, logrus.Fields{"email": email, "error": err})
					lookupErrors++
					continue
				}
				if user == nil {
					notFoundList = append(notFoundList, email)
					continue
				}

				already, err := data.Signin.HasSigninOnDate(user.ID, entType, meetingDate, span.Context())
				if err != nil {
					logging.Error(s, "Failed to check existing signin: "+email, i.Member.User, span, logrus.Fields{"email": email, "error": err})
					lookupErrors++
					continue
				}
				if already {
					alreadyList = append(alreadyList, email)
					continue
				}

				// Midnight-UTC timestamp: the same day convention DQuery in data/signin.go
				// uses, so date-based queries attribute the signin to the meeting's date.
				if _, err := data.Signin.CreateWithTimestamp(user.ID, entType, meetingDate.Truncate(24*time.Hour), span.Context()); err != nil {
					logging.Error(s, "Failed to create signin: "+email, i.Member.User, span, logrus.Fields{"email": email, "error": err})
					lookupErrors++
					continue
				}
				importedCount++
			}

			note := ""
			if !csvDeleted {
				note = csvManualDeleteNote
			}
			message := importSummary(signinTypeString, dateString, importedCount, alreadyList, notFoundList, invalidLines, lookupErrors, importMaxListItems) + note
			if len(message) >= 2000 {
				// Fall back to counts only; Discord caps message content at 2000 chars.
				message = importSummary(signinTypeString, dateString, importedCount, alreadyList, notFoundList, invalidLines, lookupErrors, 0) + note
			}

			if editErr := helpers.IntRespondEdit(s, i, message); editErr != nil {
				logging.Error(s, editErr.Error(), i.Member.User, span, logrus.Fields{"error": editErr})
				return
			}

			logging.Debug(s, fmt.Sprintf(
				"Import complete for %s on %s: %d imported, %d already signed in, %d not found, %d invalid rows, %d lookup errors",
				signinTypeString, dateString, importedCount, len(alreadyList), len(notFoundList), len(invalidLines), lookupErrors,
			), i.Member.User, span)
		}
}

// deleteCSVMessage removes the channel message carrying the uploaded CSV.
// Slash command interactions do not include the invocation message's ID, so
// recent channel history is scanned for a bot-authored message containing the
// attachment.
func deleteCSVMessage(s *discordgo.Session, i *discordgo.InteractionCreate, file *discordgo.MessageAttachment) error {
	if file == nil || file.ID == "" {
		return fmt.Errorf("no CSV attachment to delete")
	}
	if s.State == nil || s.State.User == nil {
		return fmt.Errorf("session state unavailable; cannot locate CSV message")
	}

	messages, err := s.ChannelMessages(i.ChannelID, 25, "", "", "")
	if err != nil {
		return fmt.Errorf("fetch channel messages: %w", err)
	}
	for _, message := range messages {
		if message.Author == nil || message.Author.ID != s.State.User.ID {
			continue
		}
		for _, attachment := range message.Attachments {
			if attachment.ID == file.ID {
				return s.ChannelMessageDelete(i.ChannelID, message.ID)
			}
		}
	}
	return fmt.Errorf("CSV message not found in recent channel messages")
}

// importSummary builds the final import report. maxList caps how many emails or
// row numbers are listed per section; 0 means counts only.
func importSummary(signinType, dateString string, importedCount int, alreadyList, notFoundList []string, invalidLines []int, lookupErrors, maxList int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Import complete for **%s** on %s:\n", signinType, dateString))
	if importedCount > 0 {
		fmt.Fprintf(&b, "- Imported: %d signins\n", importedCount)
	} else {
		b.WriteString("- Imported: 0 signins (nothing was added)\n")
	}
	if len(alreadyList) > 0 {
		fmt.Fprintf(&b, "- Already signed in (skipped): %d — %s\n", len(alreadyList), cappedList(alreadyList, maxList))
	}
	if len(notFoundList) > 0 {
		fmt.Fprintf(&b, "- No matching OBIII user (skipped): %d — %s\n", len(notFoundList), cappedList(notFoundList, maxList))
	}
	if len(invalidLines) > 0 {
		numbers := make([]string, len(invalidLines))
		for idx, line := range invalidLines {
			numbers[idx] = strconv.Itoa(line)
		}
		fmt.Fprintf(&b, "- Invalid rows: %d — %s\n", len(invalidLines), cappedList(numbers, maxList))
	}
	if lookupErrors > 0 {
		fmt.Fprintf(&b, "- Lookup errors: %d\n", lookupErrors)
	}
	return b.String()
}

// cappedList renders up to max items from list, appending "…and N more" when truncated.
func cappedList(list []string, max int) string {
	if len(list) <= max {
		return strings.Join(list, ", ")
	}
	return strings.Join(list[:max], ", ") + fmt.Sprintf(" …and %d more", len(list)-max)
}

// parseCsv downloads an attachment and reads its CSV records, preserving
// the first row so the caller can decide whether it contains column headers.
func parseCsv(file *discordgo.MessageAttachment) ([][]string, error) {
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
