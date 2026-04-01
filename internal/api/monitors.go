package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListMonitors(ctx context.Context, search string, tags []string, status string) ([]Monitor, error) {
	params := url.Values{}
	if search != "" {
		params.Set("name", search)
	}
	for _, tag := range tags {
		params.Add("monitor_tags", tag)
	}
	if status != "" {
		// Datadog v1 uses group_states to filter
	}

	path := "/v1/monitor"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var monitors []Monitor
	if err := json.Unmarshal(raw, &monitors); err != nil {
		return nil, err
	}

	if status != "" {
		filtered := make([]Monitor, 0)
		for _, m := range monitors {
			if m.Status == status {
				filtered = append(filtered, m)
			}
		}
		return filtered, nil
	}

	return monitors, nil
}

func (c *Client) GetMonitor(ctx context.Context, id int) (*Monitor, error) {
	path := fmt.Sprintf("/v1/monitor/%d", id)
	return doAndDecode[Monitor](c, ctx, http.MethodGet, path, nil)
}

func (c *Client) SearchMonitors(ctx context.Context, query string, status string) ([]Monitor, error) {
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}

	path := "/v1/monitor/search"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Monitors []Monitor `json:"monitors"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	if status != "" {
		filtered := make([]Monitor, 0)
		for _, m := range resp.Monitors {
			if m.Status == status {
				filtered = append(filtered, m)
			}
		}
		return filtered, nil
	}

	return resp.Monitors, nil
}

func (c *Client) MuteMonitor(ctx context.Context, id int, end string, reason string) error {
	body := map[string]any{}
	if end != "" {
		body["end"] = end
	}
	if reason != "" {
		body["scope"] = "*"
	}
	path := fmt.Sprintf("/v1/monitor/%d/mute", id)
	_, err := c.do(ctx, http.MethodPost, path, body)
	return err
}

func (c *Client) UnmuteMonitor(ctx context.Context, id int) error {
	path := fmt.Sprintf("/v1/monitor/%d/unmute", id)
	_, err := c.do(ctx, http.MethodPost, path, map[string]any{"scope": "*", "all_scopes": true})
	return err
}
