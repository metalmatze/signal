// Copyright 2020 by the contributors.
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

package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsHandler(t *testing.T) {
	reg := prometheus.NewRegistry()
	handler := NewMetricsHandler(NewHandler(), reg)

	for _, name := range []string{"aaa", "bbb", "ccc"} {
		check := func() error {
			return nil
		}

		handler.AddReadinessCheck(name, check)
		handler.AddLivenessCheck(name, check)
	}

	errorsChecks := []string{"ddd", "eee", "fff"}
	for _, name := range errorsChecks {
		check := func() error {
			return fmt.Errorf("failing health check %q", name)
		}

		handler.AddReadinessCheck(name, check)
		handler.AddLivenessCheck(name, check)
	}

	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	metricsHandler.ServeHTTP(rr, req)

	lines := strings.Split(rr.Body.String(), "\n")
	var relevantLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "healthcheck") {
			relevantLines = append(relevantLines, line)
		}
	}
	sort.Strings(relevantLines)
	actualMetrics := strings.Join(relevantLines, "\n")
	expectedMetrics := strings.TrimSpace(`
healthcheck{check="live",name="aaa"} 1
healthcheck{check="live",name="bbb"} 1
healthcheck{check="live",name="ccc"} 1
healthcheck{check="live",name="ddd"} 0
healthcheck{check="live",name="eee"} 0
healthcheck{check="live",name="fff"} 0
healthcheck{check="ready",name="aaa"} 1
healthcheck{check="ready",name="bbb"} 1
healthcheck{check="ready",name="ccc"} 1
healthcheck{check="ready",name="ddd"} 0
healthcheck{check="ready",name="eee"} 0
healthcheck{check="ready",name="fff"} 0
`)
	if actualMetrics != expectedMetrics {
		t.Errorf("expected metrics:\n%s\n\nactual metrics:\n%s\n", expectedMetrics, actualMetrics)
	}
}

func TestNewMetricsHandlerEndpoints(t *testing.T) {
	handler := NewMetricsHandler(NewHandler(), prometheus.NewRegistry())
	handler.AddReadinessCheck("fail", func() error {
		return fmt.Errorf("failing readiness check")
	})

	tests := []struct {
		name    string
		path    string
		handler http.Handler
		expect  int
	}{
		{
			name:    "default /live endpoint",
			path:    "/live",
			handler: handler,
			expect:  200,
		},
		{
			name:    "default /ready endpoint",
			path:    "/ready",
			handler: handler,
			expect:  503,
		},
		{
			name:    "custom /live endpoint",
			path:    "/",
			handler: http.HandlerFunc(handler.LiveEndpoint),
			expect:  200,
		},
		{
			name:    "custom /ready endpoint",
			path:    "/",
			handler: http.HandlerFunc(handler.ReadyEndpoint),
			expect:  503,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			tt.handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expect, rr.Code, "%s: wrong status", tt.name)
		})
	}
}
