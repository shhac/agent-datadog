package monitors

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/shhac/agent-dd/internal/api"
	"github.com/shhac/agent-dd/internal/cli/shared"
	agenterrors "github.com/shhac/agent-dd/internal/errors"
	"github.com/shhac/agent-dd/internal/output"
)

func Register(root *cobra.Command, globals func() *shared.GlobalFlags) {
	mon := &cobra.Command{
		Use:   "monitors",
		Short: "Monitor status and management",
	}

	registerList(mon, globals)
	registerGet(mon, globals)
	registerSearch(mon, globals)
	registerMute(mon, globals)
	registerUnmute(mon, globals)
	registerLLMHelp(mon)

	root.AddCommand(mon)
}

func registerList(parent *cobra.Command, globals func() *shared.GlobalFlags) {
	var search, status, tag string
	var full bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List monitors",
		RunE: func(cmd *cobra.Command, args []string) error {
			g := globals()
			return shared.WithClient(g.Org, g.Timeout, func(ctx context.Context, client *api.Client) error {
				var tags []string
				if tag != "" {
					tags = []string{tag}
				}
				monitors, err := client.ListMonitors(ctx, search, tags, status)
				if err != nil {
					return err
				}

				if full {
					shared.WritePaginatedList(shared.ToAnySlice(monitors), nil, g.Format)
					return nil
				}

				compact := make([]api.MonitorCompact, len(monitors))
				for i, m := range monitors {
					compact[i] = api.MonitorCompact{
						ID:     m.ID,
						Name:   m.Name,
						Status: m.Status,
						Type:   m.Type,
					}
				}
				shared.WritePaginatedList(shared.ToAnySlice(compact), nil, g.Format)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&search, "search", "", "Filter by name")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status (alert, warn, ok, no_data)")
	cmd.Flags().StringVar(&tag, "tag", "", "Filter by tag")
	cmd.Flags().BoolVar(&full, "full", false, "Show full monitor details")
	parent.AddCommand(cmd)
}

func registerGet(parent *cobra.Command, globals func() *shared.GlobalFlags) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get monitor details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g := globals()
			id, err := strconv.Atoi(args[0])
			if err != nil {
				output.WriteError(os.Stderr, agenterrors.Newf(agenterrors.FixableByAgent, "invalid monitor ID %q — must be an integer", args[0]))
				return nil
			}
			return shared.WithClient(g.Org, g.Timeout, func(ctx context.Context, client *api.Client) error {
				monitor, err := client.GetMonitor(ctx, id)
				if err != nil {
					return err
				}
				output.PrintJSON(monitor, true)
				return nil
			})
		},
	}
	parent.AddCommand(cmd)
}

func registerSearch(parent *cobra.Command, globals func() *shared.GlobalFlags) {
	var query, status string

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search monitors",
		RunE: func(cmd *cobra.Command, args []string) error {
			g := globals()
			if query == "" {
				output.WriteError(os.Stderr, agenterrors.New("--query is required", agenterrors.FixableByAgent))
				return nil
			}
			return shared.WithClient(g.Org, g.Timeout, func(ctx context.Context, client *api.Client) error {
				monitors, err := client.SearchMonitors(ctx, query, status)
				if err != nil {
					return err
				}
				compact := make([]api.MonitorCompact, len(monitors))
				for i, m := range monitors {
					compact[i] = api.MonitorCompact{
						ID:     m.ID,
						Name:   m.Name,
						Status: m.Status,
						Type:   m.Type,
					}
				}
				shared.WritePaginatedList(shared.ToAnySlice(compact), nil, g.Format)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "Search query (required)")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status (alert, warn, ok, no_data)")
	parent.AddCommand(cmd)
}

func registerMute(parent *cobra.Command, globals func() *shared.GlobalFlags) {
	var end, reason string

	cmd := &cobra.Command{
		Use:   "mute <id>",
		Short: "Mute a monitor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g := globals()
			id, err := strconv.Atoi(args[0])
			if err != nil {
				output.WriteError(os.Stderr, agenterrors.Newf(agenterrors.FixableByAgent, "invalid monitor ID %q", args[0]))
				return nil
			}

			var endStr string
			if end != "" {
				t, err := shared.ParseTime(end)
				if err != nil {
					output.WriteError(os.Stderr, err)
					return nil
				}
				endStr = fmt.Sprintf("%d", t.Unix())
			}

			return shared.WithClient(g.Org, g.Timeout, func(ctx context.Context, client *api.Client) error {
				if err := client.MuteMonitor(ctx, id, endStr, reason); err != nil {
					return err
				}
				output.PrintJSON(map[string]any{
					"status":     "muted",
					"monitor_id": id,
				}, true)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&end, "end", "", "Mute end time (relative or absolute)")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for muting")
	parent.AddCommand(cmd)
}

func registerUnmute(parent *cobra.Command, globals func() *shared.GlobalFlags) {
	cmd := &cobra.Command{
		Use:   "unmute <id>",
		Short: "Unmute a monitor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g := globals()
			id, err := strconv.Atoi(args[0])
			if err != nil {
				output.WriteError(os.Stderr, agenterrors.Newf(agenterrors.FixableByAgent, "invalid monitor ID %q", args[0]))
				return nil
			}
			return shared.WithClient(g.Org, g.Timeout, func(ctx context.Context, client *api.Client) error {
				if err := client.UnmuteMonitor(ctx, id); err != nil {
					return err
				}
				output.PrintJSON(map[string]any{
					"status":     "unmuted",
					"monitor_id": id,
				}, true)
				return nil
			})
		},
	}
	parent.AddCommand(cmd)
}
