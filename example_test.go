package signal

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/metalmatze/signal/log"
	"github.com/metalmatze/signal/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ExampleSignals() {
	logger := log.NewLogftmLogger(os.Stdout)

	requests := prometheus.NewCounterVec(prom.CounterOpts{
		Name: "requests_total",
		Help: "Counting the requests",
	}, []attribute.Key{"code", "method"})
	reg := prom.NewRegistry()
	reg.MustRegister(requests)

	tracer := otel.Tracer("main")

	successHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode := 200
		w.WriteHeader(statusCode)

		// define attributes once
		message := attribute.String("msg", "request handled")
		method := attribute.String("method", r.Method)
		code := attribute.Int("code", statusCode)

		// reuse attributes across metrics, logs, traces
		requests.With(method, code).Inc()
		logger.Log(method, code, message)
		_, span := tracer.Start(r.Context(), "successHandler", trace.WithAttributes(method, code, message))
		defer span.End()
	})
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode := 500
		w.WriteHeader(statusCode)

		// define attributes once
		message := attribute.String("msg", "request handled")
		method := attribute.String("method", r.Method)
		code := attribute.Int("code", statusCode)

		// reuse attributes across metrics, logs, traces
		requests.With(method, code).Inc()
		logger.Log(method, code, message)
		_, span := tracer.Start(r.Context(), "successHandler", trace.WithAttributes(method, code, message))
		defer span.End()
	})

	sendRequest(successHandler, http.MethodGet)
	sendRequest(errorHandler, http.MethodPut)
	sendRequest(successHandler, http.MethodGet)
	sendRequest(errorHandler, http.MethodGet)

	// Output:
	// method=GET code=200 msg="request handled"
	// method=PUT code=500 msg="request handled"
	// method=GET code=200 msg="request handled"
	// method=GET code=500 msg="request handled"
}

func sendRequest(handler http.Handler, method string) {
	req := httptest.NewRequest(method, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}
