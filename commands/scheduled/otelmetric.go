package scheduled

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/ritsec/ops-bot-iii/config"
)

var (
	obiiiOTLPOnce      sync.Once
	obiiiBotProvider   *sdkmetric.MeterProvider
	obiiiBotLiveness   metric.Int64Counter
)

// obiiiOTLPInterval returns the OTLP export cadence from commands.heartbeat.otlp.interval
// (default 60s). The cadence is aligned to the down-window detection window (~30-60s), NOT the
// 10s liveness probe, so datapoint volume stays low while absence still trips within one window.
func obiiiOTLPInterval() time.Duration {
	const def = "60s"
	raw := config.GetString("commands.heartbeat.otlp.interval")
	if raw == "" {
		raw = def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		return 60 * time.Second
	}
	return d
}

// InitBotOTLP registers the OTel metric pipeline once: an OTLP HTTP exporter pointed at Datadog's
// agentless intake endpoint, a PeriodicReader that flushes every obiiiOTLPInterval, and an explicit
// MeterProvider (no process-global). Safe to call from multiple scheduled tasks. It is a no-op when
// no datadog.api_key is configured (e.g. kuma-only deployments).
func InitBotOTLP() {
	obiiiOTLPOnce.Do(func() {
		if ddAPIKey == "" {
			logrus.Warn("commands.scheduled.otelmetric: no commands.heartbeat.datadog.api_key set; OTLP metric export disabled")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		site := ddSite
		if site == "" {
			site = "datadoghq.com"
		}

		exporter, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint("api."+site), // host only, no scheme/path; TLS (https) default
			otlpmetrichttp.WithURLPath("/otlp/v1/metrics"), // Datadog agentless OTLP metrics intake path
			otlpmetrichttp.WithHeaders(map[string]string{
				"DD-API-KEY": ddAPIKey,
				// Convert resource attributes (service.name, env) into per-series tags so staging/prod are distinguishable.
				"dd-otel-metric-config": `{"resource_attributes_as_tags": true}`,
			}),
			otlpmetrichttp.WithTemporalitySelector(deltaForCounters),
		)
		if err != nil {
			logrus.Errorf("commands.scheduled.otelmetric: create OTLP exporter: %v", err)
			return
		}

		reader := sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(obiiiOTLPInterval()),
		)

		provider := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(botResource()),
			sdkmetric.WithReader(reader),
		)

		meter := provider.Meter("OBIII")
		c, err := meter.Int64Counter(
			"obiii.bot",
			metric.WithDescription("OBIII bot liveness heartbeat; 1 per successful tick"),
		)
		if err != nil {
			logrus.Errorf("commands.scheduled.otelmetric: create counter: %v", err)
			return
		}

		obiiiBotProvider = provider
		obiiiBotLiveness = c
		logrus.Infof("commands.scheduled.otelmetric: OTLP metric export enabled (interval=%s, site=%s)", obiiiOTLPInterval(), site)
	})
}

// deltaForCounters forces Delta temporality for counters so each periodic export carries only the
// increment since the last export (low cost). Alive => points flow every interval; dead => the series
// stops after at most one window, so a Datadog On Missing Data monitor trips within ~1 interval.
func deltaForCounters(ik sdkmetric.InstrumentKind) metricdata.Temporality {
	if ik == sdkmetric.InstrumentKindCounter {
		return metricdata.DeltaTemporality
	}
	return metricdata.CumulativeTemporality
}

// botResource builds the OTel resource so service.name + env become tags on every emitted series.
func botResource() *resource.Resource {
	name := config.GetString("name")
	if name == "" {
		name = "OBIII"
	}
	attrs := []attribute.KeyValue{attribute.String("service.name", name)}
	if env := config.GetString("env"); env != "" {
		attrs = append(attrs, attribute.String("env", env))
	}
	return resource.NewSchemaless(attrs...)
}

// RecordBotLiveness increments the liveness counter. Best-effort: a failure must never break the
// heartbeat loop.
func RecordBotLiveness() {
	if obiiiBotLiveness != nil {
		obiiiBotLiveness.Add(context.Background(), 1)
	}
}

// ShutdownBotOTLP flushes pending metrics and stops the periodic reader on graceful stop.
func ShutdownBotOTLP() {
	if obiiiBotProvider == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := obiiiBotProvider.Shutdown(ctx); err != nil {
		logrus.Warnf("commands.scheduled.otelmetric: shutdown meter provider: %v", err)
	}
}
