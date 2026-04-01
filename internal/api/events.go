package api

import (
	"context"
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
	resp, err := doAndDecode[EventListResponse](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Events, nil
}

func (c *Client) GetEvent(ctx context.Context, id int64) (*Event, error) {
	type eventResp struct {
		Event Event `json:"event"`
	}
	path := fmt.Sprintf("/v1/events/%d", id)
	return doAndDecodeField[eventResp, Event](c, ctx, http.MethodGet, path, nil, func(r *eventResp) *Event { return &r.Event })
}
