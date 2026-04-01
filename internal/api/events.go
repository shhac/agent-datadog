package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EventListResponse struct {
	Events []Event `json:"events"`
}

func (c *Client) ListEvents(ctx context.Context, from, to int64, source string, tags []string) ([]Event, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", from)},
		"end":   {fmt.Sprintf("%d", to)},
	}
	if source != "" {
		params.Set("sources", source)
	}
	for _, tag := range tags {
		params.Add("tags", tag)
	}

	path := "/v1/events?" + params.Encode()
	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp EventListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return resp.Events, nil
}

func (c *Client) GetEvent(ctx context.Context, id int64) (*Event, error) {
	path := fmt.Sprintf("/v1/events/%d", id)

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Event Event `json:"event"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp.Event, nil
}
