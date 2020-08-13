package signal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/metalmatze/signal"
	"github.com/metalmatze/signal/healthcheck"
	"github.com/metalmatze/signal/internalserver"
	"github.com/prometheus/client_golang/prometheus"
)

func Example() {
	registry := prometheus.NewRegistry()
	healthchecks := healthcheck.NewMetricsHandler(healthcheck.NewHandler(), registry)

	h := internalserver.NewHandler(
		internalserver.WithHealthchecks(healthchecks),
		internalserver.WithPrometheusRegistry(registry),
		internalserver.WithPProf(),
	)

	fmt.Print(dumpRequest(h, http.MethodGet, "/"))

	// Output:
	// HTTP/1.1 200 OK
	// Connection: close
	// Content-Type: text/html; charset=utf-8
	//
	// <html><head><title>Internal</title></head><body>
	// <p><a href='/debug'>/debug - Exposes pprof endpoints to consume via HTTP</a></p>
	// <p><a href='/live'>/live - Exposes liveness checks</a></p>
	// <p><a href='/metrics'>/metrics - Exposes Prometheus metrics</a></p>
	// <p><a href='/ready'>/ready - Exposes readiness checks</a></p>
	// </body></html>
}

func dumpRequest(handler http.Handler, method string, path string) string {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	dump, err := httputil.DumpResponse(rr.Result(), true)
	if err != nil {
		panic(err)
	}
	return strings.Replace(string(dump), "\r\n", "\n", -1)
}

func Example_signal() {
	signals := signal.New()
	signals.Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)) // overwriting to get rid of timestamps for example

	event := signals.Event(
		signal.Log("msg", "event happened"),
		signal.Counter("counter", "we count events"),
	)

	counting1 := func(i int) func(e *signal.Event) {
		return func(e *signal.Event) {
			e.Log("i", i)
		}
	}
	counting2 := func(e *signal.Event) func(i int) {
		return func(i int) {
			e.Log("i", i)
		}
	}

	for i := 0; i < 10; i++ {
		counting1(i)(event)
		counting2(event)(i)

		event.Happens(func() {
		})
	}

	// Output:
	// msg="event happened" i=0
	// msg="event happened" i=1
	// msg="event happened" i=2
	// msg="event happened" i=3
	// msg="event happened" i=4
	// msg="event happened" i=5
	// msg="event happened" i=6
	// msg="event happened" i=7
	// msg="event happened" i=8
	// msg="event happened" i=9
}
func dumpSignals(signals *signal.Signal) string {
	return ""
}
