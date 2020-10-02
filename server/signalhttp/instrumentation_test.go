package signalhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func Example() {
	registry := prometheus.NewRegistry()

	// Create the HTTP handler instrumenter with extra labels group and handler.
	// The labels code and method are added out-of-the-box.
	instrumenter := NewHandlerInstrumenter(registry, []string{"group", "handler"})

	r := http.NewServeMux()

	r.Handle("/", instrumenter.NewHandler(
		prometheus.Labels{"group": "apiv1", "handler": "index"},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("success"))
		}),
	))

	// Make a HTTP request against the internalserver Handler
	fmt.Print(dumpRequest(r, http.MethodGet, "/"))

	// Output:
	// HTTP/1.1 200 OK
	// Connection: close
	//
	// success
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
