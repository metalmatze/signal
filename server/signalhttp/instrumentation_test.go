// Copyright 2021 by the contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
