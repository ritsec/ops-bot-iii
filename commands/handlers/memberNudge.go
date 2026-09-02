package handlers

import (
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// memberNudgeRoles are the role IDs that mark a user as verified, matching
	// the roles the /member flow can assign
	memberNudgeRoles = []string{
		config.GetString("commands.member.member_role_id"),
		config.GetString("commands.member.external_role_id"),
		config.GetString("commands.member.prospective_role_id"),
		config.GetString("commands.member.staff_role_id"),
		config.GetString("commands.member.alumni_role_id"),
	}

	// nudgeHelpChannel is the channel name where nagging stops and where
	// nudging never happens
	nudgeHelpChannel = "bot-help"

	// nudgeBackoffBase is the wait before the first re-nudge
	nudgeBackoffBase = 30 * time.Second

	// nudgeBackoffCap is the maximum wait between re-nudges
	nudgeBackoffCap = 30 * time.Minute

	// nudgeDeletionDelay is how long a nudge and its triggering message stay up
	nudgeDeletionDelay = 1 * time.Minute
)

// nudgeState tracks the exponential-backoff position for one user
type nudgeState struct {
	lastNudge time.Time
	backoff   time.Duration
}

// nudges tracks per-user nudge backoff state, keyed by user ID
var nudges = struct {
	sync.Mutex
	m map[string]*nudgeState
}{m: make(map[string]*nudgeState)}

func nudgeClear(userID string) {
	nudges.Lock()
	delete(nudges.m, userID)
	nudges.Unlock()
}

// nudgeDue reports whether the user may be nudged now and records the nudge
// with the next backoff step if so
func nudgeDue(userID string) bool {
	nudges.Lock()
	defer nudges.Unlock()

	now := time.Now()
	st, ok := nudges.m[userID]
	if !ok {
		// first nudge is immediate
		nudges.m[userID] = &nudgeState{
			lastNudge: now,
			backoff:   nudgeBackoffBase,
		}
		return true
	}

	if time.Since(st.lastNudge) < st.backoff {
		return false
	}

	next := st.backoff * 2
	if next > nudgeBackoffCap {
		next = nudgeBackoffCap
	}
	st.lastNudge = now
	st.backoff = next
	return true
}

// nudgeVerified reports whether the user holds any assignable member role
func nudgeVerified(m *discordgo.MessageCreate) bool {
	if m.Member == nil {
		return false
	}
	for _, roleID := range memberNudgeRoles {
		if roleID == "" {
			continue
		}
		for _, role := range m.Member.Roles {
			if role == roleID {
				return true
			}
		}
	}
	return false
}

// MemberNudge nudges unverified users to run /member. It is a MessageCreate
// handler: when a non-bot user who has no assignable member role posts in any
// channel other than bot-help, it replies publicly with a pointer to /member.
// Successive nudges back off exponentially (30s, 1m, 2m, ... capped at 30m).
// One minute after each nudge, both the bot's reply and the triggering user
// message are deleted. Nagging stops once the user is verified or posts in
// bot-help.
func MemberNudge(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	if m.GuildID == "" {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	// never nudge in bot-help; a post there ends the nagging
	if strings.EqualFold(channel.Name, nudgeHelpChannel) {
		nudgeClear(m.Author.ID)
		return
	}
	// cheap in-memory gate first: skip the data-layer lookup entirely while a
	// prior nudge is still within its backoff window
	if !nudgeDue(m.Author.ID) {
		return
	}

	span := tracer.StartSpan(
		"commands.handlers.memberNudge:MemberNudge",
		tracer.ResourceName("Handlers.MemberNudge"),
	)
	defer span.Finish()

	// verified via an assignable role or the data layer: stop nagging
	if nudgeVerified(m) || data.User.IsVerified(m.Author.ID, span.Context()) {
		nudgeClear(m.Author.ID)
		return
	}

	message, err := s.ChannelMessageSendReply(m.ChannelID, helpers.AtUser(m.Author.ID)+" — you are not a member. Please run `/member` to verify.", m.Reference())
	if err != nil {
		logging.Error(s, err.Error(), m.Author, span, logrus.Fields{"error": err})
		return
	}

	logging.Debug(s, "Nudged unverified user to /member", m.Author, span)

	// delete the reply and the triggering message one minute after the nudge;
	// the goroutine outlives the request span, so it runs under its own child
	go func(spanCtx ddtrace.SpanContext) {
		delSpan := tracer.StartSpan(
			"commands.handlers.memberNudge:MemberNudge.delete",
			tracer.ResourceName("Handlers.MemberNudge.delete"),
			tracer.ChildOf(spanCtx),
		)
		defer delSpan.Finish()

		time.Sleep(nudgeDeletionDelay)
		if delErr := s.ChannelMessageDelete(m.ChannelID, message.ID); delErr != nil {
			logging.Error(s, delErr.Error(), m.Author, delSpan, logrus.Fields{"error": delErr})
		}
		if delErr := s.ChannelMessageDelete(m.ChannelID, m.ID); delErr != nil {
			logging.Error(s, delErr.Error(), m.Author, delSpan, logrus.Fields{"error": delErr})
		}
	}(span.Context())
}
