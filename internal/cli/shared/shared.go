package shared

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shhac/agent-dd/internal/api"
	"github.com/shhac/agent-dd/internal/config"
	"github.com/shhac/agent-dd/internal/credential"
	agenterrors "github.com/shhac/agent-dd/internal/errors"
	"github.com/shhac/agent-dd/internal/output"
)

type GlobalFlags struct {
	Org     string
	Format  string
	Timeout int
}

func MakeContext(timeoutMs int) (context.Context, context.CancelFunc) {
	if timeoutMs > 0 {
		return context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	}
	return context.WithCancel(context.Background())
}

func ResolveOrg(orgAlias string) (string, error) {
	if orgAlias != "" {
		return orgAlias, nil
	}
	if env := os.Getenv("DD_ORG"); env != "" {
		return env, nil
	}
	cfg := config.Read()
	if cfg.DefaultOrg != "" {
		return cfg.DefaultOrg, nil
	}
	available := make([]string, 0)
	for name := range cfg.Organizations {
		available = append(available, name)
	}
	hint := "No organizations configured. Add one with 'agent-dd org add <alias>'"
	if len(available) > 0 {
		hint = fmt.Sprintf("Available organizations: %s. Set a default with 'agent-dd org set-default <alias>'", strings.Join(available, ", "))
	}
	return "", agenterrors.New("no organization specified", agenterrors.FixableByAgent).WithHint(hint)
}

func NewClientFromOrg(orgAlias string) (*api.Client, error) {
	// Check env vars first (standard DD env vars)
	if orgAlias == "" {
		apiKey := os.Getenv("DD_API_KEY")
		appKey := os.Getenv("DD_APP_KEY")
		if apiKey != "" && appKey != "" {
			site := os.Getenv("DD_SITE")
			if site == "" {
				site = "datadoghq.com"
			}
			return api.NewClient(apiKey, appKey, site), nil
		}
	}

	alias, err := ResolveOrg(orgAlias)
	if err != nil {
		return nil, err
	}

	cred, err := credential.Get(alias)
	if err != nil {
		var nf *credential.NotFoundError
		if errors.As(err, &nf) {
			return nil, agenterrors.Newf(agenterrors.FixableByHuman, "credentials for organization %q not found", alias).
				WithHint("Add credentials with 'agent-dd org add " + alias + " --api-key <key> --app-key <key>'")
		}
		return nil, agenterrors.Wrap(err, agenterrors.FixableByHuman)
	}

	if cred.APIKey == "" {
		return nil, agenterrors.Newf(agenterrors.FixableByHuman, "organization %q has no API key", alias).
			WithHint("Update with 'agent-dd org update " + alias + " --api-key <key>'")
	}

	cfg := config.Read()
	site := "datadoghq.com"
	if org, ok := cfg.Organizations[alias]; ok && org.Site != "" {
		site = org.Site
	}

	return api.NewClient(cred.APIKey, cred.AppKey, site), nil
}

var ClientFactory func() (*api.Client, error)

func WithClient(orgAlias string, timeout int, fn func(ctx context.Context, client *api.Client) error) error {
	ctx, cancel := MakeContext(timeout)
	defer cancel()

	var client *api.Client
	var err error
	if ClientFactory != nil {
		client, err = ClientFactory()
	} else {
		client, err = NewClientFromOrg(orgAlias)
	}
	if err != nil {
		output.WriteError(os.Stderr, err)
		return nil
	}

	if err := fn(ctx, client); err != nil {
		output.WriteError(os.Stderr, err)
	}
	return nil
}

func ToAnySlice[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

func WritePaginatedList(items []any, pagination *output.Pagination, format string) {
	f := output.ResolveFormat(format)
	if f == output.FormatNDJSON {
		w := output.NewNDJSONWriter(os.Stdout)
		for _, item := range items {
			w.WriteItem(item)
		}
		if pagination != nil {
			w.WritePagination(pagination)
		}
		return
	}
	result := map[string]any{"data": items}
	if pagination != nil {
		result["pagination"] = pagination
	}
	output.PrintJSON(result, true)
}

// ParseTime parses relative (now-15m), RFC3339, or unix epoch time strings.
func ParseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	if s == "now" {
		return time.Now(), nil
	}

	// Relative: now-15m, now-1h, now-1d, now-7d, now+1h
	if strings.HasPrefix(s, "now") {
		return parseRelativeTime(s)
	}

	// RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Unix epoch seconds
	if epoch, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(epoch, 0), nil
	}

	return time.Time{}, agenterrors.Newf(agenterrors.FixableByAgent,
		"invalid time format %q — use relative (now-15m), RFC3339 (2024-01-15T10:00:00Z), or unix epoch", s)
}

func parseRelativeTime(s string) (time.Time, error) {
	now := time.Now()
	rest := s[3:] // strip "now"

	if rest == "" {
		return now, nil
	}

	var sign time.Duration = -1
	if rest[0] == '+' {
		sign = 1
		rest = rest[1:]
	} else if rest[0] == '-' {
		rest = rest[1:]
	} else {
		return time.Time{}, agenterrors.Newf(agenterrors.FixableByAgent, "invalid relative time %q", s)
	}

	if len(rest) < 2 {
		return time.Time{}, agenterrors.Newf(agenterrors.FixableByAgent, "invalid relative time %q", s)
	}

	unit := rest[len(rest)-1]
	numStr := rest[:len(rest)-1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return time.Time{}, agenterrors.Newf(agenterrors.FixableByAgent, "invalid relative time %q", s)
	}

	var duration time.Duration
	switch unit {
	case 's':
		duration = time.Duration(num) * time.Second
	case 'm':
		duration = time.Duration(num) * time.Minute
	case 'h':
		duration = time.Duration(num) * time.Hour
	case 'd':
		duration = time.Duration(num) * 24 * time.Hour
	case 'w':
		duration = time.Duration(num) * 7 * 24 * time.Hour
	default:
		return time.Time{}, agenterrors.Newf(agenterrors.FixableByAgent,
			"invalid time unit %q in %q — use s, m, h, d, or w", string(unit), s)
	}

	return now.Add(sign * duration), nil
}

// ParseTimeDefaultFrom returns the parsed --from time, defaulting to 1 hour ago.
func ParseTimeDefaultFrom(s string) (time.Time, error) {
	if s == "" {
		return time.Now().Add(-1 * time.Hour), nil
	}
	return ParseTime(s)
}

// ParseTimeDefaultTo returns the parsed --to time, defaulting to now.
func ParseTimeDefaultTo(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	return ParseTime(s)
}
