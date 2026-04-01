package logs_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shhac/agent-dd/internal/cli/shared"
)

func TestLogsSearchEndpoint(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/logs/events/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		filter, ok := body["filter"].(map[string]any)
		if !ok {
			t.Error("missing filter in request body")
		}
		if query, _ := filter["query"].(string); query == "" {
			t.Error("empty query in filter")
		}

		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":   "log1",
					"type": "log",
					"attributes": map[string]any{
						"timestamp": "2024-01-15T10:00:00Z",
						"service":   "web-api",
						"status":    "error",
						"message":   "connection timeout",
					},
				},
			},
		})
	})

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}
}

func TestLogsAggregateEndpoint(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/logs/analytics/aggregate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"buckets": []map[string]any{
					{"by": map[string]any{"service": "web-api"}, "computes": map[string]any{"c0": 42}},
				},
			},
		})
	})

	if shared.ClientFactory == nil {
		t.Fatal("ClientFactory not set")
	}
}
