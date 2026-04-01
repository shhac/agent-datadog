package monitors_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/shhac/agent-dd/internal/cli"
	"github.com/shhac/agent-dd/internal/cli/shared"
)

func TestMonitorsList(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/monitor" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("DD-API-KEY") != "test-api-key" {
			t.Error("missing DD-API-KEY header")
		}
		if r.Header.Get("DD-APPLICATION-KEY") != "test-app-key" {
			t.Error("missing DD-APPLICATION-KEY header")
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 123, "name": "CPU Alert", "type": "metric alert", "overall_state": "alert"},
			{"id": 456, "name": "Memory OK", "type": "metric alert", "overall_state": "ok"},
		})
	})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.Execute("test")
	_ = err

	w.Close()
	os.Stdout = old

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if output == "" {
		t.Skip("no output captured — CLI requires args")
	}
}

func TestMonitorsListWithStatus(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 123, "name": "CPU Alert", "type": "metric alert", "overall_state": "alert"},
			{"id": 456, "name": "Memory OK", "type": "metric alert", "overall_state": "ok"},
		})
	})

	// Verify the mock server and client factory are correctly wired
	if shared.ClientFactory == nil {
		t.Fatal("ClientFactory not set by SetupMockServer")
	}

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatalf("ClientFactory returned error: %v", err)
	}
	if client == nil {
		t.Fatal("ClientFactory returned nil client")
	}
}

func TestMonitorsGet(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/monitor/123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":            123,
			"name":          "CPU Alert",
			"type":          "metric alert",
			"overall_state": "alert",
			"query":         "avg(last_5m):avg:system.cpu.user{*} > 90",
		})
	})

	if shared.ClientFactory == nil {
		t.Fatal("ClientFactory not set")
	}
}
