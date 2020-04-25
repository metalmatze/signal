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
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTCPDialCheck(t *testing.T) {
	assert.NoError(t, TCPDialCheck("google.com:80", 5*time.Second)())
	assert.Error(t, TCPDialCheck("google.com:25327", 5*time.Second)())
}

func TestHTTPGetCheck(t *testing.T) {
	assert.NoError(t, HTTPGetCheck("https://www.google.com", 5*time.Second)())
	assert.Error(t, HTTPGetCheck("http://google.com", 5*time.Second)(), "redirect should fail")
	assert.Error(t, HTTPGetCheck("https://google.com/nonexistent", 5*time.Second)(), "404 should fail")
}

func TestHTTPCheck(t *testing.T) {
	assert.NoError(t, HTTPCheck("https://www.google.com", http.MethodGet, http.StatusOK, 5*time.Second)())
	assert.NoError(t, HTTPCheck("https://www.google.com", http.MethodPost, http.StatusMethodNotAllowed, 5*time.Second)())
	assert.Error(t, HTTPCheck("https://www.google.com", http.MethodPost, http.StatusOK, 5*time.Second)())
	assert.NoError(t, HTTPCheck("https://www.google.com/nonexistent", http.MethodGet, http.StatusNotFound, 5*time.Second)())
	assert.Error(t, HTTPCheck("https://www.google.com/nonexistent", http.MethodGet, http.StatusOK, 5*time.Second)())
}

func TestDatabasePingCheck(t *testing.T) {
	assert.Error(t, DatabasePingCheck(nil, 1*time.Second)(), "nil DB should fail")

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	assert.NoError(t, DatabasePingCheck(db, 1*time.Second)(), "ping should succeed")
}

func TestDNSResolveCheck(t *testing.T) {
	assert.NoError(t, DNSResolveCheck("google.com", 5*time.Second)())
	assert.Error(t, DNSResolveCheck("nonexistent.google.com", 5*time.Second)())
}

func TestGoroutineCountCheck(t *testing.T) {
	assert.NoError(t, GoroutineCountCheck(1000)())
	assert.Error(t, GoroutineCountCheck(0)())
}
