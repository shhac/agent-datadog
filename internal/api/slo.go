package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SLOListResponse struct {
	Data []SLO `json:"data"`
}

func (c *Client) ListSLOs(ctx context.Context, search string, tags []string) ([]SLO, error) {
	params := url.Values{}
	if search != "" {
		params.Set("query", search)
	}
	if len(tags) > 0 {
		for _, tag := range tags {
			params.Add("tags_query", tag)
		}
	}

	path := "/v1/slo"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp SLOListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) GetSLO(ctx context.Context, id string) (*SLO, error) {
	path := fmt.Sprintf("/v1/slo/%s", url.PathEscape(id))

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data SLO `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) GetSLOHistory(ctx context.Context, id string, from, to int64) (*SLOHistory, error) {
	params := url.Values{
		"from_ts": {fmt.Sprintf("%d", from)},
		"to_ts":   {fmt.Sprintf("%d", to)},
	}

	path := fmt.Sprintf("/v1/slo/%s/history?%s", url.PathEscape(id), params.Encode())

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data SLOHistory `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
