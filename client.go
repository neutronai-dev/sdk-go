package neutron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	nucleusID  string
	tenantID   string
	httpClient *http.Client
}

type Options struct {
	BaseURL    string
	Token      string
	NucleusID  string
	TenantID   string
	HTTPClient *http.Client
}

type APIError struct {
	Status int
	Body   []byte
}

func (err APIError) Error() string {
	return fmt.Sprintf("Neutron AI request failed: %d", err.Status)
}

func NewClient(options Options) (*Client, error) {
	if strings.TrimSpace(options.BaseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if options.NucleusID != "" && options.TenantID != "" && options.NucleusID != options.TenantID {
		return nil, fmt.Errorf("nucleus ID and tenant ID must match when both are provided")
	}
	httpClient := options.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL:    strings.TrimRight(options.BaseURL, "/"),
		token:      options.Token,
		nucleusID:  firstNonEmpty(options.NucleusID, options.TenantID),
		tenantID:   options.TenantID,
		httpClient: httpClient,
	}, nil
}

func (client *Client) request(ctx context.Context, method string, path string, body map[string]any, output any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	request, err := http.NewRequestWithContext(ctx, method, client.baseURL+path, reader)
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/json")
	request.Header.Set("user-agent", "neutron-ai-go/0.3.0")
	if client.token != "" {
		request.Header.Set("authorization", "Bearer "+client.token)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return APIError{Status: response.StatusCode, Body: rawBody}
	}
	if output == nil || len(rawBody) == 0 {
		return nil
	}
	return json.Unmarshal(rawBody, output)
}

func (client *Client) withNucleus(body map[string]any) (map[string]any, error) {
	output := map[string]any{}
	for key, value := range body {
		output[key] = value
	}
	nucleusID, _ := output["nucleusId"].(string)
	if nucleusID == "" {
		nucleusID, _ = output["nucleus_id"].(string)
	}
	tenantID, _ := output["tenantId"].(string)
	if tenantID == "" {
		tenantID, _ = output["tenant_id"].(string)
	}
	if nucleusID != "" && tenantID != "" && nucleusID != tenantID {
		return nil, fmt.Errorf("nucleus ID and tenant ID must match when both are provided")
	}
	nucleusID = firstNonEmpty(nucleusID, tenantID, client.nucleusID, client.tenantID)
	if strings.TrimSpace(nucleusID) == "" {
		return nil, fmt.Errorf("nucleus ID is required")
	}
	output["nucleusId"] = nucleusID
	output["tenantId"] = nucleusID
	delete(output, "nucleus_id")
	delete(output, "tenant_id")
	return output, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (client *Client) Health(ctx context.Context) (map[string]any, error) {
	output := map[string]any{}
	err := client.request(ctx, http.MethodGet, "/health", nil, &output)
	return output, err
}

func (client *Client) Remember(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/remember", scopedBody, &output)
	return output, err
}

func (client *Client) Recall(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/recall", scopedBody, &output)
	return output, err
}

func (client *Client) AgentContext(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/agent-context", scopedBody, &output)
	return output, err
}

func (client *Client) Reflect(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/reflect", scopedBody, &output)
	return output, err
}

func (client *Client) LearnSkill(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/skills/learn", scopedBody, &output)
	return output, err
}

func (client *Client) UpsertKnowledge(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/kb/upsert", scopedBody, &output)
	return output, err
}

func (client *Client) Forget(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/forget", scopedBody, &output)
	return output, err
}

func (client *Client) Compact(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/compact", scopedBody, &output)
	return output, err
}

func (client *Client) Audit(ctx context.Context, scopeID string, limit int, since int64, until int64) (map[string]any, error) {
	if client.nucleusID == "" {
		return nil, fmt.Errorf("nucleus ID is required")
	}
	if limit <= 0 {
		limit = 50
	}
	query := url.Values{
		"nucleusId": []string{client.nucleusID},
		"scopeId":   []string{scopeID},
		"limit":     []string{fmt.Sprintf("%d", limit)},
	}
	if since > 0 {
		query.Set("since", fmt.Sprintf("%d", since))
	}
	if until > 0 {
		query.Set("until", fmt.Sprintf("%d", until))
	}
	output := map[string]any{}
	err := client.request(ctx, http.MethodGet, "/v1/memory/audit?"+query.Encode(), nil, &output)
	return output, err
}

func (client *Client) CacheStats(ctx context.Context, agentID string) (map[string]any, error) {
	if client.nucleusID == "" {
		return nil, fmt.Errorf("nucleus ID is required")
	}
	query := url.Values{"nucleusId": []string{client.nucleusID}}
	if agentID != "" {
		query.Set("agentId", agentID)
	}
	output := map[string]any{}
	err := client.request(ctx, http.MethodGet, "/v1/memory/cache/stats?"+query.Encode(), nil, &output)
	return output, err
}

func (client *Client) InvalidateCache(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/cache/invalidate", scopedBody, &output)
	return output, err
}

func (client *Client) WarmCache(ctx context.Context, body map[string]any) (map[string]any, error) {
	scopedBody, err := client.withNucleus(body)
	if err != nil {
		return nil, err
	}
	output := map[string]any{}
	err = client.request(ctx, http.MethodPost, "/v1/memory/cache/warm", scopedBody, &output)
	return output, err
}

func (client *Client) ListMCPTools(ctx context.Context) (map[string]any, error) {
	output := map[string]any{}
	err := client.request(ctx, http.MethodGet, "/mcp/tools", nil, &output)
	return output, err
}
