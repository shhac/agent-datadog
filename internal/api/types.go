package api

import "encoding/json"

// Monitor represents a Datadog monitor.
type Monitor struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Query    string           `json:"query,omitempty"`
	Message  string           `json:"message,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
	Status   string           `json:"overall_state,omitempty"`
	Created  string           `json:"created,omitempty"`
	Modified string           `json:"modified,omitempty"`
	Options  *json.RawMessage `json:"options,omitempty"`
}

// MonitorCompact is the token-efficient view of a monitor.
type MonitorCompact struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// LogEntry represents a Datadog log entry.
type LogEntry struct {
	ID         string            `json:"id"`
	Timestamp  string            `json:"timestamp,omitempty"`
	Service    string            `json:"service,omitempty"`
	Status     string            `json:"status,omitempty"`
	Message    string            `json:"message,omitempty"`
	Host       string            `json:"host,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	Attributes map[string]any    `json:"attributes,omitempty"`
}

// LogEntryCompact is the token-efficient view of a log entry.
type LogEntryCompact struct {
	Timestamp string `json:"timestamp"`
	Service   string `json:"service,omitempty"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// MetricSeries represents a metric query result series.
type MetricSeries struct {
	Metric string      `json:"metric,omitempty"`
	Tags   []string    `json:"tags,omitempty"`
	Points [][]float64 `json:"points"`
}

// MetricMetadata represents metadata about a metric.
type MetricMetadata struct {
	Name        string `json:"metric,omitempty"`
	Type        string `json:"type,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
	Integration string `json:"integration,omitempty"`
	PerUnit     string `json:"per_unit,omitempty"`
	ShortName   string `json:"short_name,omitempty"`
}

// Event represents a Datadog event.
type Event struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	Text       string   `json:"text,omitempty"`
	DateHappened int64  `json:"date_happened,omitempty"`
	Source     string   `json:"source,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Priority   string   `json:"priority,omitempty"`
	AlertType  string   `json:"alert_type,omitempty"`
	Host       string   `json:"host,omitempty"`
}

// Host represents a Datadog host.
type Host struct {
	Name       string   `json:"name"`
	Aliases    []string `json:"aliases,omitempty"`
	Apps       []string `json:"apps,omitempty"`
	IsMuted    bool     `json:"is_muted"`
	MuteTimeout int64   `json:"mute_timeout,omitempty"`
	Sources    []string `json:"sources,omitempty"`
	Up         bool     `json:"up"`
	TagsBySource map[string][]string `json:"tags_by_source,omitempty"`
	LastReportedTime int64 `json:"last_reported_time,omitempty"`
}

// Incident represents a Datadog incident.
type Incident struct {
	ID         string              `json:"id"`
	Type       string              `json:"type,omitempty"`
	Attributes *IncidentAttributes `json:"attributes,omitempty"`
}

type IncidentAttributes struct {
	Title         string `json:"title,omitempty"`
	Status        string `json:"status,omitempty"`
	Severity      string `json:"severity,omitempty"`
	Created       string `json:"created,omitempty"`
	Modified      string `json:"modified,omitempty"`
	CommanderUser *IncidentUser `json:"commander_user,omitempty"`
}

type IncidentUser struct {
	Handle string `json:"handle,omitempty"`
	Name   string `json:"name,omitempty"`
}

// Trace span from APM.
type TraceSpan struct {
	TraceID  string  `json:"trace_id"`
	SpanID   string  `json:"span_id"`
	Service  string  `json:"service,omitempty"`
	Name     string  `json:"name,omitempty"`
	Resource string  `json:"resource,omitempty"`
	Type     string  `json:"type,omitempty"`
	Start    int64   `json:"start,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Error    int     `json:"error,omitempty"`
	Status   string  `json:"status,omitempty"`
}

// APMService represents an APM service.
type APMService struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

// SLO represents a Service Level Objective.
type SLO struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Thresholds  []SLOThreshold `json:"thresholds,omitempty"`
	Status      *SLOStatus  `json:"overall_status,omitempty"`
}

type SLOThreshold struct {
	Timeframe string  `json:"timeframe"`
	Target    float64 `json:"target"`
	Warning   float64 `json:"warning,omitempty"`
}

type SLOStatus struct {
	Status    float64 `json:"status,omitempty"`
	ErrorBudgetRemaining float64 `json:"error_budget_remaining,omitempty"`
}

// SLOHistory represents SLO history data.
type SLOHistory struct {
	Overall  *SLOHistoryMetrics   `json:"overall,omitempty"`
	Thresholds map[string]SLOHistoryMetrics `json:"thresholds,omitempty"`
}

type SLOHistoryMetrics struct {
	SLIValue          float64 `json:"sli_value,omitempty"`
	SpanPrecision     float64 `json:"span_precision,omitempty"`
	Uptime            float64 `json:"uptime,omitempty"`
	ErrorBudgetRemaining float64 `json:"error_budget_remaining,omitempty"`
}
