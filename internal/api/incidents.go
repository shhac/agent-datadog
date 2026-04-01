package api

import (
	"context"
	"net/http"
	"net/url"
)

type IncidentListResponse struct {
	Data []Incident `json:"data"`
}

func (c *Client) ListIncidents(ctx context.Context, status string) ([]Incident, error) {
	params := url.Values{}
	if status != "" {
		params.Set("filter[status]", status)
	}

	path := "/v2/incidents"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}

	resp, err := doAndDecode[IncidentListResponse](c, ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) GetIncident(ctx context.Context, id string) (*Incident, error) {
	type resp struct {
		Data Incident `json:"data"`
	}
	return doAndDecodeField[resp, Incident](c, ctx, http.MethodGet, "/v2/incidents/"+url.PathEscape(id), nil, func(r *resp) *Incident { return &r.Data })
}

func (c *Client) CreateIncident(ctx context.Context, title, severity, commanderHandle string) (*Incident, error) {
	data := map[string]any{
		"type": "incidents",
		"attributes": map[string]any{
			"title": title,
			"fields": map[string]any{
				"severity": map[string]any{
					"type":  "dropdown",
					"value": severity,
				},
			},
		},
	}

	if commanderHandle != "" {
		data["relationships"] = map[string]any{
			"commander_user": map[string]any{
				"data": map[string]any{
					"type": "users",
					"id":   commanderHandle,
				},
			},
		}
	}

	type resp struct {
		Data Incident `json:"data"`
	}
	return doAndDecodeField[resp, Incident](c, ctx, http.MethodPost, "/v2/incidents", map[string]any{"data": data}, func(r *resp) *Incident { return &r.Data })
}

func (c *Client) UpdateIncident(ctx context.Context, id string, status, severity string) (*Incident, error) {
	attrs := map[string]any{}
	if status != "" {
		attrs["status"] = status
	}
	if severity != "" {
		attrs["fields"] = map[string]any{
			"severity": map[string]any{
				"type":  "dropdown",
				"value": severity,
			},
		}
	}

	body := map[string]any{
		"data": map[string]any{
			"type":       "incidents",
			"id":         id,
			"attributes": attrs,
		},
	}

	type resp struct {
		Data Incident `json:"data"`
	}
	return doAndDecodeField[resp, Incident](c, ctx, http.MethodPatch, "/v2/incidents/"+url.PathEscape(id), body, func(r *resp) *Incident { return &r.Data })
}
