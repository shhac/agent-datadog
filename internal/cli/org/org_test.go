package org_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/shhac/agent-dd/internal/cli/shared"
)

func TestOrgTestCommand(t *testing.T) {
	shared.SetupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/validate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("DD-API-KEY") == "" {
			t.Error("missing DD-API-KEY")
		}
		if r.Header.Get("DD-APPLICATION-KEY") == "" {
			t.Error("missing DD-APPLICATION-KEY")
		}
		json.NewEncoder(w).Encode(map[string]any{"valid": true})
	})

	client, err := shared.ClientFactory()
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("nil client")
	}
}
