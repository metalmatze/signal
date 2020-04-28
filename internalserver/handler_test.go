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
