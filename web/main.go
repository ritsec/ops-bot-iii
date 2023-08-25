package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
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

	err := Router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
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
	err := Router.Run(config.Web.Hostname + ":" + config.Web.Port)
	if err != nil {
		logging.CriticalDD("Error running web server", span, logrus.Fields{"error": err})
	}
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
