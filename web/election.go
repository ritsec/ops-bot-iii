package web

import (
	"github.com/gin-gonic/gin"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// vote is the handler for the /vote/:id endpoint
func vote(ctx *gin.Context) {
	span, _ := tracer.StartSpanFromContext(
		ctx.Request.Context(),
		"web.election:vote",
		tracer.ResourceName("/vote/:id"),
	)
	defer span.Finish()

	voteID := ctx.Param("id")

	entVoteResult, err := data.VoteResult.Get(voteID, span.Context())
	if err != nil {
		logging.ErrorDD("Failed to get vote result", span, logrus.Fields{"error": err})
		ctx.String(500, err.Error())
		return
	}

	ctx.Data(200, "text/html", []byte(entVoteResult.HTML))
}
