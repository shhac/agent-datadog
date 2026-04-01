package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shhac/agent-dd/internal/api"
)

func SetupMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(func() {
		srv.Close()
		ClientFactory = nil
	})
	ClientFactory = func() (*api.Client, error) {
		return api.NewTestClient(srv.URL, "test-api-key", "test-app-key"), nil
	}
	return srv
}
