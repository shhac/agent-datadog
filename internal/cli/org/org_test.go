package org_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shhac/agent-dd/internal/cli/shared"
)

func TestOrgValidate(t *testing.T) {
	var called bool
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.URL.Path != "/api/v1/validate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("DD-API-KEY") != "test-api-key" {
			t.Error("missing or wrong DD-API-KEY")
		}
		if r.Header.Get("DD-APPLICATION-KEY") != "test-app-key" {
			t.Error("missing or wrong DD-APPLICATION-KEY")
		}
		json.NewEncoder(w).Encode(map[string]any{"valid": true})
	})

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Validate(context.Background()); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if !called {
		t.Error("mock handler was never called")
	}
}
