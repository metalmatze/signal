package signal

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/metalmatze/signal/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Signal struct {
	Logger       log.Logger
	Registry     *prometheus.Registry
	Healthchecks healthcheck.Handler
	// TODO: Add Tracing immediately
}

func New() *Signal {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	healthchecks := healthcheck.NewMetricsHandler(healthcheck.NewHandler(), reg)

	{
		jaeger, _ := jaeger.NewRawExporter(jaeger.WithAgentEndpoint(""))
		tracer, _ := trace.NewProvider(
			trace.WithConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}),
			trace.WithSyncer(jaeger),
		)
		global.SetTraceProvider(tracer)
	}

	return &Signal{
		Logger:       logger,
		Registry:     reg,
		Healthchecks: healthchecks,
	}
}

type EventOption func(s *Signal, e *Event)

func Log(keyvals ...interface{}) EventOption {
	return func(s *Signal, e *Event) {
		e.logger = log.With(s.Logger, keyvals...)
	}
}

func Counter(name, help string) EventOption {
	return func(s *Signal, e *Event) {
		c := prometheus.NewCounter(prometheus.CounterOpts{Name: name, Help: help})
		s.Registry.MustRegister(c)
	}
}

func (s *Signal) Event(option ...EventOption) *Event {
	e := &Event{
		signal: s,
	}

	for _, o := range option {
		o(s, e)
	}

	return e
}

type Event struct {
	signal *Signal
	logger log.Logger
}

func (e *Event) Log(keyvals ...interface{}) *Event {
	if e.logger != nil {
		_ = e.logger.Log(keyvals...)
	}
	return e
}

func (e *Event) Count(count float64) {

}

func (e *Event) Happens(f func()) {
	f()
}
