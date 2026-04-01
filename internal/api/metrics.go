package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MetricQueryResponse struct {
	Status string         `json:"status"`
	Series []MetricSeries `json:"series"`
}

func (c *Client) QueryMetrics(ctx context.Context, query string, from, to int64) (*MetricQueryResponse, error) {
	params := url.Values{
		"query": {query},
		"from":  {fmt.Sprintf("%d", from)},
		"to":    {fmt.Sprintf("%d", to)},
	}

	path := "/v1/query?" + params.Encode()
	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp MetricQueryResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type MetricListResponse struct {
	Data []MetricListEntry `json:"data"`
}

type MetricListEntry struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (c *Client) ListMetrics(ctx context.Context, search string, tag string) (*MetricListResponse, error) {
	params := url.Values{}
	if search != "" {
		params.Set("filter[configured]", "true")
		params.Set("filter[tags_configured]", search)
	}
	if tag != "" {
		params.Set("filter[tags]", tag)
	}

	// v2 metrics listing can be limited; fall back to v1 search
	if search != "" {
		return c.searchMetricsV1(ctx, search)
	}

	path := "/v2/metrics"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp MetricListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) searchMetricsV1(ctx context.Context, query string) (*MetricListResponse, error) {
	params := url.Values{"q": {"metrics:" + query}}
	path := "/v1/search?" + params.Encode()

	raw, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Results struct {
			Metrics []string `json:"metrics"`
		} `json:"results"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	entries := make([]MetricListEntry, len(resp.Results.Metrics))
	for i, m := range resp.Results.Metrics {
		entries[i] = MetricListEntry{ID: m, Type: "metric"}
	}
	return &MetricListResponse{Data: entries}, nil
}

func (c *Client) GetMetricMetadata(ctx context.Context, metricName string) (*MetricMetadata, error) {
	path := fmt.Sprintf("/v1/metrics/%s", url.PathEscape(metricName))
	return doAndDecode[MetricMetadata](c, ctx, http.MethodGet, path, nil)
}
