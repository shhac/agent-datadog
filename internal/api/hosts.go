package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type HostListResponse struct {
	HostList   []Host `json:"host_list"`
	TotalReturned int `json:"total_returned"`
	TotalMatching int `json:"total_matching"`
}

func (c *Client) ListHosts(ctx context.Context, search string, tags []string) (*HostListResponse, error) {
	params := url.Values{}
	if search != "" {
		params.Set("filter", search)
	}
	if len(tags) > 0 {
		for _, tag := range tags {
			params.Add("filter", tag)
		}
	}

	path := "/v1/hosts"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp HostListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetHost(ctx context.Context, hostname string) (*Host, error) {
	params := url.Values{"filter": {hostname}}
	path := "/v1/hosts?" + params.Encode()

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp HostListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if len(resp.HostList) == 0 {
		return nil, fmt.Errorf("host %q not found", hostname)
	}
	return &resp.HostList[0], nil
}

func (c *Client) MuteHost(ctx context.Context, hostname string, end int64, reason string) error {
	body := map[string]any{
		"hostname": hostname,
	}
	if end > 0 {
		body["end"] = end
	}
	if reason != "" {
		body["message"] = reason
	}

	path := fmt.Sprintf("/v1/host/%s/mute", url.PathEscape(hostname))
	_, err := c.do(ctx, http.MethodPost, path, body)
	return err
}
