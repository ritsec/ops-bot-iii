package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// Router is the gin router
	Router *gin.Engine
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()

	Router.SetTrustedProxies(nil)
}

// Start starts the web server
func Start(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"web.main:Start",
		tracer.ResourceName("web.main:Start"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	loadRoutes(span.Context())
	Router.Run(config.Web.Hostname + ":" + config.Web.Port)
}

// loadRoutes loads the routes
func loadRoutes(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"web.main:loadRoutes",
		tracer.ResourceName("web.main:loadRoutes"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	Router.GET("/vote/:id", vote)
}
