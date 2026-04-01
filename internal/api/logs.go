package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type LogSearchRequest struct {
	Filter *LogFilter `json:"filter"`
	Sort   string     `json:"sort,omitempty"`
	Page   *LogPage   `json:"page,omitempty"`
}

type LogFilter struct {
	Query string `json:"query"`
	From  string `json:"from"`
	To    string `json:"to"`
}

type LogPage struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type LogSearchResponse struct {
	Data []LogData      `json:"data"`
	Meta *LogSearchMeta `json:"meta,omitempty"`
}

type LogData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes LogAttributes  `json:"attributes"`
}

type LogAttributes struct {
	Timestamp  string         `json:"timestamp,omitempty"`
	Service    string         `json:"service,omitempty"`
	Status     string         `json:"status,omitempty"`
	Message    string         `json:"message,omitempty"`
	Host       string         `json:"host,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type LogSearchMeta struct {
	Page *LogSearchMetaPage `json:"page,omitempty"`
}

type LogSearchMetaPage struct {
	After string `json:"after,omitempty"`
}

func (c *Client) SearchLogs(ctx context.Context, query, from, to, sort string, limit int) (*LogSearchResponse, error) {
	req := LogSearchRequest{
		Filter: &LogFilter{
			Query: query,
			From:  from,
			To:    to,
		},
	}
	if sort != "" {
		req.Sort = sort
	}
	if limit > 0 {
		req.Page = &LogPage{Limit: limit}
	}

	raw, err := c.do(ctx, http.MethodPost, "/v2/logs/events/search", req)
	if err != nil {
		return nil, err
	}

	var resp LogSearchResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type LogAggregateBucket struct {
	Computes map[string]any `json:"computes,omitempty"`
	By       map[string]any `json:"by,omitempty"`
}

type LogAggregateResponse struct {
	Data struct {
		Buckets []LogAggregateBucket `json:"buckets"`
	} `json:"data"`
}

func (c *Client) AggregateLogs(ctx context.Context, query, from, to string, groupBy []string) (*LogAggregateResponse, error) {
	computes := []map[string]any{
		{"aggregation": "count", "type": "total"},
	}

	groups := make([]map[string]any, 0, len(groupBy))
	for _, g := range groupBy {
		groups = append(groups, map[string]any{
			"facet": g,
			"limit": 10,
			"total": map[string]string{"aggregation": "count", "order": "desc"},
		})
	}

	body := map[string]any{
		"filter": map[string]any{
			"query": query,
			"from":  from,
			"to":    to,
		},
		"compute": computes,
	}
	if len(groups) > 0 {
		body["group_by"] = groups
	}

	raw, err := c.do(ctx, http.MethodPost, "/v2/logs/analytics/aggregate", body)
	if err != nil {
		return nil, err
	}

	var resp LogAggregateResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
