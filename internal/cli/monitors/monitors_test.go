package monitors_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

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

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatal(err)
	}

	monitors, err := client.ListMonitors(context.Background(), "", nil, "")
	if err != nil {
		t.Fatalf("ListMonitors failed: %v", err)
	}
	if len(monitors) != 2 {
		t.Fatalf("expected 2 monitors, got %d", len(monitors))
	}
	if monitors[0].Name != "CPU Alert" {
		t.Errorf("expected name 'CPU Alert', got %q", monitors[0].Name)
	}
}

func TestMonitorsListWithStatus(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 123, "name": "CPU Alert", "type": "metric alert", "overall_state": "alert"},
			{"id": 456, "name": "Memory OK", "type": "metric alert", "overall_state": "ok"},
		})
	})

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatal(err)
	}

	monitors, err := client.ListMonitors(context.Background(), "", nil, "alert")
	if err != nil {
		t.Fatalf("ListMonitors failed: %v", err)
	}
	if len(monitors) != 1 {
		t.Fatalf("expected 1 filtered monitor, got %d", len(monitors))
	}
	if monitors[0].Status != "alert" {
		t.Errorf("expected status 'alert', got %q", monitors[0].Status)
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

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatal(err)
	}

	monitor, err := client.GetMonitor(context.Background(), 123)
	if err != nil {
		t.Fatalf("GetMonitor failed: %v", err)
	}
	if monitor.Name != "CPU Alert" {
		t.Errorf("expected name 'CPU Alert', got %q", monitor.Name)
	}
	if monitor.Query != "avg(last_5m):avg:system.cpu.user{*} > 90" {
		t.Errorf("unexpected query: %q", monitor.Query)
	}
}
