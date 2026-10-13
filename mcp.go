package neutron

import "strings"

func CreateNeutronMCPServerURL(baseURL string, path string) string {
	if path == "" {
		path = "/mcp"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimRight(baseURL, "/") + path
}

func CreateOpenAIMCPTool(serverURL string, serverLabel string) map[string]any {
	if serverLabel == "" {
		serverLabel = "neutron_memory"
	}
	return map[string]any{
		"type":             "mcp",
		"server_label":     serverLabel,
		"server_url":       serverURL,
		"require_approval": "never",
	}
}

func CreateAnthropicMCPConnector(serverURL string, serverLabel string, authorization string) map[string]any {
	if serverLabel == "" {
		serverLabel = "neutron_memory"
	}
	server := map[string]any{
		"type": "url",
		"url":  serverURL,
		"name": serverLabel,
	}
	if authorization != "" {
		server["authorization_token"] = authorization
	}
	return map[string]any{
		"mcp_servers": []map[string]any{server},
		"tools": []map[string]any{{
			"type":            "mcp_toolset",
			"mcp_server_name": serverLabel,
		}},
		"betas": []string{"mcp-client-2025-11-20"},
	}
}
