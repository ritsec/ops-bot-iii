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

// nudgeState tracks the per-user backoff position and the guild where the
// last timeout was applied
type nudgeState struct {
	lastNudge time.Time
	backoff   time.Duration
	guildID   string
}

// nudges tracks per-user nudge backoff state, keyed by user ID
var nudges = struct {
	sync.Mutex
	m map[string]*nudgeState
}{m: make(map[string]*nudgeState)}

// nudgeClear removes the user's nudge state and lifts any active Discord
// timeout; lifting is best-effort so a missing MODERATE_MEMBERS permission
// does not break the nagging flow
func nudgeClear(s *discordgo.Session, userID string, user *discordgo.User, span ddtrace.Span) {
	nudges.Lock()
	st, ok := nudges.m[userID]
	delete(nudges.m, userID)
	guildID := ""
	if ok {
		guildID = st.guildID
	}
	nudges.Unlock()

	if guildID == "" {
		return
	}

	if err := s.GuildMemberTimeout(guildID, userID, nil); err != nil {
		logging.Error(s, "Failed to clear guild member timeout", user, span, logrus.Fields{"error": err})
	} else {
		logging.Debug(s, "Cleared guild member timeout", user, span)
	}
}

// nudgeDue reports whether the user may be nudged now and records the nudge
// with the next backoff step if so; the returned duration is the Discord
// member timeout to apply for this nudge
func nudgeDue(userID, guildID string) (bool, time.Duration) {
	nudges.Lock()
	defer nudges.Unlock()

	now := time.Now()
	st, ok := nudges.m[userID]
	if !ok {
		// first nudge is immediate
		nudges.m[userID] = &nudgeState{
			lastNudge: now,
			backoff:   nudgeBackoffBase,
			guildID:   guildID,
		}
		return true, nudgeBackoffBase
	}

	if time.Since(st.lastNudge) < st.backoff {
		return false, 0
	}

	next := st.backoff * 2
	if next > nudgeBackoffCap {
		next = nudgeBackoffCap
	}
	st.lastNudge = now
	st.backoff = next
	st.guildID = guildID
	return true, next
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
// channel other than bot-help, it replies publicly with a pointer to /member
// and times the user out (communication_disabled_until) for the current
// backoff step, which doubles with every nudge (30s, 1m, 2m, ... capped at
// 30m), so the user cannot post again until the silence has lifted. One
// minute after each nudge, both the bot's reply and the triggering user
// message are deleted. Nagging stops once the user is verified or posts in
// bot-help, which also clears any active timeout. Applying the timeout
// requires MODERATE_MEMBERS permission; when it fails, the in-memory backoff
// gate continues to throttle re-nudges as a fallback.
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

	span := tracer.StartSpan(
		"commands.handlers.memberNudge:MemberNudge",
		tracer.ResourceName("Handlers.MemberNudge"),
	)
	defer span.Finish()

	// never nudge in bot-help; a post there ends the nagging and lifts any
	// active timeout
	if strings.EqualFold(channel.Name, nudgeHelpChannel) {
		nudgeClear(s, m.Author.ID, m.Author, span)
		return
	}

	// in-memory backoff gate first; also the fallback throttle when the
	// Discord-native timeout below cannot be applied
	due, timeout := nudgeDue(m.Author.ID, m.GuildID)
	if !due {
		return
	}

	// verified via an assignable role or the data layer: stop nagging
	if nudgeVerified(m) || data.User.IsVerified(m.Author.ID, span.Context()) {
		nudgeClear(s, m.Author.ID, m.Author, span)
		return
	}

	message, err := s.ChannelMessageSendReply(m.ChannelID, helpers.AtUser(m.Author.ID)+" — you are not a member. Please run `/member` to verify.", m.Reference())
	if err != nil {
		logging.Error(s, err.Error(), m.Author, span, logrus.Fields{"error": err})
		return
	}

	logging.Debug(s, "Nudged unverified user to /member", m.Author, span)

	// silence the user for the current backoff step so they cannot post (and
	// re-trigger the nudge) until the window has elapsed; a failure here is
	// non-fatal — the in-memory gate above keeps throttling
	until := time.Now().Add(timeout)
	if timeoutErr := s.GuildMemberTimeout(m.GuildID, m.Author.ID, &until); timeoutErr != nil {
		logging.Error(s, "Failed to apply guild member timeout", m.Author, span, logrus.Fields{"error": timeoutErr})
	}

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
