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

package internalserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_AddEndpoint(t *testing.T) {
	h := NewHandler()
	h.AddEndpoint("/foo", "some handler", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("bar"))
	})

	assert.Len(t, h.endpoints, 1)
	assert.Equal(t, endpoint{Pattern: "/foo", Description: "some handler"}, h.endpoints[0])

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/foo", nil)
	h.ServeHTTP(rr, req)

	assert.Equal(t, "bar", rr.Body.String())
}
