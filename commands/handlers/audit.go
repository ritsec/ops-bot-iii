package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/helpers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// MemberJoin is a handler for when a member joins the server
func MemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	span := tracer.StartSpan(
		"commands.handlers.audit:MemberJoin",
		tracer.ResourceName("Handlers.MemberJoin"),
	)
	defer span.Finish()

	logging.AuditButton(
		s,
		fmt.Sprintf("%v Joined", helpers.AtUser(m.Member.User.ID)),
		discordgo.Button{
			Label: "View Profile",
			URL:   fmt.Sprintf("https://discordapp.com/users/%v/", m.Member.User.ID),
			Style: discordgo.LinkButton,
		},
		m.Member.User,
		span,
	)

	err := helpers.SendDirectMessage(s, m.Member.User.ID, "Welcome to the RITSEC Discord Server! Please read `#readme` channel and to be verified as a member use `/member`", span.Context())
	if err != nil {
		logging.Error(s, err.Error(), m.Member.User, span, logrus.Fields{"error": err})
	}
}

// MemberJoin is a handler for when a member joins the server
func MemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	span := tracer.StartSpan(
		"commands.handlers.audit:MemberLeave",
		tracer.ResourceName("Handlers.MemberLeave"),
	)
	defer span.Finish()

	logging.AuditButton(
		s,
		fmt.Sprintf("%v Left", helpers.AtUser(m.Member.User.ID)),
		discordgo.Button{
			Label: "View Profile",
			URL:   fmt.Sprintf("https://discordapp.com/users/%v/", m.Member.User.ID),
			Style: discordgo.LinkButton,
		},
		m.Member.User,
		span,
	)
}
