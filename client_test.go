package neutron

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRememberInjectsNucleusLegacyAliasAndToken(t *testing.T) {
	var requestBody map[string]any
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader = request.Header.Get("authorization")
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		_ = json.NewEncoder(writer).Encode(map[string]any{"memoryId": "mem_1"})
	}))
	defer server.Close()

	client, err := NewClient(Options{BaseURL: server.URL, Token: "secret", NucleusID: "nucleus_demo"})
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Remember(context.Background(), map[string]any{
		"scopeId": "user:derek",
		"type":    "preference",
		"text":    "hello",
	})
	if err != nil {
		t.Fatal(err)
	}

	if authHeader != "Bearer secret" {
		t.Fatalf("unexpected auth header: %s", authHeader)
	}
	if requestBody["nucleusId"] != "nucleus_demo" || requestBody["tenantId"] != "nucleus_demo" {
		t.Fatalf("nucleus aliases were not injected: %#v", requestBody)
	}
	if response["memoryId"] != "mem_1" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestNewClientRejectsMismatchedNucleusAliases(t *testing.T) {
	_, err := NewClient(Options{
		BaseURL:   "https://api.example.test",
		NucleusID: "nucleus_a",
		TenantID:  "nucleus_b",
	})
	if err == nil {
		t.Fatal("expected mismatched aliases to fail")
	}
}

func TestCreateNeutronMCPServerURL(t *testing.T) {
	got := CreateNeutronMCPServerURL("https://mcp.neutronai.dev/", "mcp")
	if got != "https://mcp.neutronai.dev/mcp" {
		t.Fatalf("unexpected MCP URL: %s", got)
	}
}
