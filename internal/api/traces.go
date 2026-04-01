package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type TraceSearchResponse struct {
	Data []TraceData `json:"data"`
}

type TraceData struct {
	Type       string          `json:"type"`
	Attributes TraceAttributes `json:"attributes"`
}

type TraceAttributes struct {
	TraceID    string  `json:"trace_id,omitempty"`
	SpanID     string  `json:"span_id,omitempty"`
	Service    string  `json:"service,omitempty"`
	Name       string  `json:"name,omitempty"`
	Resource   string  `json:"resource,omitempty"`
	Type       string  `json:"type,omitempty"`
	Start      int64   `json:"start,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
	Error      int     `json:"error,omitempty"`
	Status     string  `json:"status,omitempty"`
}

func (c *Client) SearchTraces(ctx context.Context, query, service, from, to string, limit int) (*TraceSearchResponse, error) {
	filterQuery := query
	if service != "" && query == "" {
		filterQuery = "service:" + service
	} else if service != "" {
		filterQuery = "service:" + service + " " + query
	}

	body := map[string]any{
		"filter": map[string]any{
			"query": filterQuery,
			"from":  from,
			"to":    to,
		},
	}
	if limit > 0 {
		body["page"] = map[string]any{"limit": limit}
	}

	raw, err := c.do(ctx, http.MethodPost, "/v2/spans/events/search", body)
	if err != nil {
		return nil, err
	}

	var resp TraceSearchResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ServiceListResponse struct {
	Data []ServiceData `json:"data"`
}

type ServiceData struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes ServiceAttributes `json:"attributes"`
}

type ServiceAttributes struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func (c *Client) ListServices(ctx context.Context, search string) ([]APMService, error) {
	params := url.Values{}
	if search != "" {
		params.Set("filter", search)
	}

	path := "/v1/services"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// v1 services endpoint returns a map
	var serviceMap map[string]any
	if err := json.Unmarshal(raw, &serviceMap); err != nil {
		return nil, err
	}

	services := make([]APMService, 0)
	for name := range serviceMap {
		if search == "" || containsSubstring(name, search) {
			services = append(services, APMService{Name: name})
		}
	}
	return services, nil
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexSubstring(s, sub) >= 0)
}

func indexSubstring(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
