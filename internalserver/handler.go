package internalserver

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"sort"

	"github.com/metalmatze/signal/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	http.ServeMux
	endpoints []endpoint
}

type endpoint struct {
	Pattern     string
	Description string
}

func NewHandler(options ...Option) *Handler {
	h := &Handler{}

	h.HandleFunc("/", h.index)

	for _, option := range options {
		option(h)
	}

	return h
}

func (h *Handler) AddEndpoint(pattern string, description string, handler http.HandlerFunc) {
	h.endpoints = append(h.endpoints, endpoint{
		Pattern:     pattern,
		Description: description,
	})

	// Sort endpoints by pattern after adding a new one, to always show them in the same order.
	sort.Slice(h.endpoints, func(i, j int) bool {
		return h.endpoints[i].Pattern < h.endpoints[j].Pattern
	})

	h.HandleFunc(pattern, handler)
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	html := "<html><head><title>Internal</title></head><body>\n"

	for _, e := range h.endpoints {
		html += fmt.Sprintf("<p><a href='%s'>%s - %s</a></p>\n", e.Pattern, e.Pattern, e.Description)
	}
	html += `</body></html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

type Option func(h *Handler)

func WithHealthchecks(healthchecks healthcheck.Handler) Option {
	return func(h *Handler) {
		h.AddEndpoint(
			"/live",
			"Exposes liveness checks",
			healthchecks.LiveEndpoint,
		)
		h.AddEndpoint(
			"/ready",
			"Exposes readiness checks",
			healthchecks.ReadyEndpoint,
		)
	}
}

func WithPrometheusRegistry(registry *prometheus.Registry) Option {
	return func(h *Handler) {
		h.AddEndpoint(
			"/metrics",
			"Exposes Prometheus metrics",
			promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP,
		)
	}
}

func WithPProf() Option {
	return func(h *Handler) {
		m := http.NewServeMux()
		m.HandleFunc("/debug/pprof/*", pprof.Index)
		m.HandleFunc("/pprof/cmdline", pprof.Cmdline)
		m.HandleFunc("/debug/pprof/profile", pprof.Profile)
		m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		m.HandleFunc("/debug/pprof/trace", pprof.Trace)

		h.AddEndpoint(
			"/debug",
			"Exposes pprof endpoints to consume via HTTP",
			m.ServeHTTP,
		)
	}
}
