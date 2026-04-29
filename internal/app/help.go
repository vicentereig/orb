package app

import (
	"errors"

	"github.com/vicentereig/orb/internal/output"
)

type HelpDocument struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Usage       string        `json:"usage"`
	Topics      []string      `json:"topics"`
	GlobalFlags []HelpFlag    `json:"global_flags"`
	Commands    []HelpCommand `json:"commands"`
	Notes       []string      `json:"notes,omitempty"`
}

type HelpTopic struct {
	Topic    string        `json:"topic"`
	Commands []HelpCommand `json:"commands"`
	Notes    []string      `json:"notes,omitempty"`
}

type HelpCommand struct {
	Name          string     `json:"name"`
	Summary       string     `json:"summary"`
	RequiredFlags []string   `json:"required_flags,omitempty"`
	OptionalFlags []HelpFlag `json:"optional_flags,omitempty"`
	Examples      []string   `json:"examples,omitempty"`
	Notes         []string   `json:"notes,omitempty"`
}

type HelpFlag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *App) Help(topic string) string {
	if topic == "" {
		return output.Success(helpDocument(), meta("system", "help"))
	}
	help, ok := helpTopics()[topic]
	if !ok {
		return output.Error(errors.New("unknown help topic: "+topic), "usage_error", meta("system", "help"))
	}
	return output.Success(help, meta("system", "help"))
}

func helpDocument() HelpDocument {
	topics := []string{"customers", "subscriptions", "plans", "prices", "metrics", "invoices", "credit-notes", "events", "alerts", "examples"}
	commands := []HelpCommand{
		{Name: "version", Summary: "Print CLI version information.", Examples: []string{"orb version"}},
		{Name: "ping", Summary: "Check Orb API connectivity and credentials.", Examples: []string{"orb ping"}},
		{Name: "help", Summary: "Return structured command help for humans and coding agents.", OptionalFlags: []HelpFlag{{Name: "--json", Description: "Accepted for explicit agent workflows; output is JSON by default."}}, Examples: []string{"orb help", "orb help events", "orb help --json"}},
	}
	for _, topic := range topics {
		if topic == "examples" {
			continue
		}
		label := topicLabel(topic)
		commands = append(commands, HelpCommand{
			Name:     topic,
			Summary:  "Inspect " + label + ". Run `orb help " + topic + "` for command-specific flags and examples.",
			Examples: []string{"orb help " + topic},
		})
	}
	return HelpDocument{
		Name:        "orb",
		Description: "Orb billing forensics from the terminal.",
		Usage:       "orb <command> [options]",
		Topics:      topics,
		GlobalFlags: globalHelpFlags(),
		Commands:    commands,
		Notes: []string{
			"All commands return a stable JSON envelope with success, data, meta, and error fields.",
			"Read commands require ORB_API_KEY unless they are version or help commands.",
			"Use --pretty for formatted JSON; compact JSON is the default for scripts and agents.",
		},
	}
}

func helpTopics() map[string]HelpTopic {
	return map[string]HelpTopic{
		"customers": {
			Topic: "customers",
			Commands: []HelpCommand{
				{Name: "customers list", Summary: "List customers.", OptionalFlags: flags("--created-after", "--created-before", "--limit", "--cursor"), Examples: []string{"orb customers list --limit 50", "orb customers list --created-after 2026-04-01T00:00:00Z"}},
				{Name: "customers get", Summary: "Fetch one customer by Orb ID or external ID.", RequiredFlags: []string{"--id or --external-id"}, Examples: []string{"orb customers get --id cus_123", "orb customers get --external-id workspace_123"}},
				{Name: "customers costs", Summary: "Fetch customer-level costs for a timeframe.", RequiredFlags: []string{"--id or --external-id", "--from", "--to"}, OptionalFlags: flags("--currency", "--view-mode"), Examples: []string{"orb customers costs --id cus_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z"}},
				{Name: "customers credits", Summary: "List customer credits.", RequiredFlags: []string{"--id or --external-id"}, OptionalFlags: flags("--currency", "--include-all-blocks", "--limit", "--cursor"), Examples: []string{"orb customers credits --id cus_123 --include-all-blocks"}},
				{Name: "customers credit-ledger", Summary: "List customer credit ledger entries.", RequiredFlags: []string{"--id or --external-id"}, OptionalFlags: flags("--currency", "--entry-type", "--entry-status", "--limit", "--cursor"), Examples: []string{"orb customers credit-ledger --id cus_123 --entry-status committed"}},
			},
		},
		"subscriptions": {
			Topic: "subscriptions",
			Commands: []HelpCommand{
				{Name: "subscriptions list", Summary: "List subscriptions.", OptionalFlags: flags("--customer-id", "--external-customer-id", "--plan-id", "--external-plan-id", "--status", "--created-after", "--created-before", "--limit", "--cursor"), Examples: []string{"orb subscriptions list --customer-id cus_123 --status active"}},
				{Name: "subscriptions get", Summary: "Fetch one subscription.", RequiredFlags: []string{"--id"}, Examples: []string{"orb subscriptions get --id sub_123"}},
				{Name: "subscriptions usage", Summary: "Fetch subscription usage for a timeframe.", RequiredFlags: []string{"--id", "--from", "--to"}, OptionalFlags: flags("--granularity", "--group-by", "--billable-metric-id", "--first-dimension-key", "--first-dimension-value", "--second-dimension-key", "--second-dimension-value", "--view-mode"), Examples: []string{"orb subscriptions usage --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z --granularity day"}},
				{Name: "subscriptions costs", Summary: "Fetch subscription costs for a timeframe.", RequiredFlags: []string{"--id", "--from", "--to"}, OptionalFlags: flags("--currency", "--view-mode"), Examples: []string{"orb subscriptions costs --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z"}},
				{Name: "subscriptions schedule", Summary: "Fetch subscription schedule phases.", RequiredFlags: []string{"--id"}, OptionalFlags: flags("--start-after", "--start-before", "--limit", "--cursor"), Examples: []string{"orb subscriptions schedule --id sub_123"}},
			},
		},
		"plans": {
			Topic: "plans",
			Commands: []HelpCommand{
				{Name: "plans list", Summary: "List plans.", OptionalFlags: flags("--status", "--created-after", "--created-before", "--limit", "--cursor"), Examples: []string{"orb plans list --status active"}},
				{Name: "plans get", Summary: "Fetch one plan by Orb ID or external ID.", RequiredFlags: []string{"--id or --external-id"}, Examples: []string{"orb plans get --id plan_123", "orb plans get --external-id pro_plan"}},
			},
		},
		"prices": {
			Topic: "prices",
			Commands: []HelpCommand{
				{Name: "prices list", Summary: "List prices.", OptionalFlags: flags("--limit", "--cursor"), Examples: []string{"orb prices list --limit 100"}},
				{Name: "prices get", Summary: "Fetch one price by Orb ID or external ID.", RequiredFlags: []string{"--id or --external-id"}, Examples: []string{"orb prices get --id price_123", "orb prices get --external-id api_calls_price"}},
			},
		},
		"metrics": {
			Topic: "metrics",
			Commands: []HelpCommand{
				{Name: "metrics list", Summary: "List billable metrics.", OptionalFlags: flags("--created-after", "--created-before", "--limit", "--cursor"), Examples: []string{"orb metrics list"}},
				{Name: "metrics get", Summary: "Fetch one billable metric.", RequiredFlags: []string{"--id"}, Examples: []string{"orb metrics get --id metric_123"}},
			},
		},
		"invoices": {
			Topic: "invoices",
			Commands: []HelpCommand{
				{Name: "invoices list", Summary: "List invoices.", OptionalFlags: flags("--customer-id", "--external-customer-id", "--subscription-id", "--status", "--invoice-date-after", "--invoice-date-before", "--limit", "--cursor"), Examples: []string{"orb invoices list --customer-id cus_123 --status issued"}},
				{Name: "invoices get", Summary: "Fetch one invoice.", RequiredFlags: []string{"--id"}, Examples: []string{"orb invoices get --id inv_123"}},
				{Name: "invoices summary", Summary: "List invoice summaries using the same filters as invoices list.", OptionalFlags: flags("--customer-id", "--external-customer-id", "--subscription-id", "--status", "--invoice-date-after", "--invoice-date-before", "--limit", "--cursor"), Examples: []string{"orb invoices summary --customer-id cus_123"}},
				{Name: "invoices upcoming", Summary: "Fetch an upcoming invoice for a subscription.", RequiredFlags: []string{"--subscription-id"}, Examples: []string{"orb invoices upcoming --subscription-id sub_123"}},
			},
		},
		"credit-notes": {
			Topic: "credit-notes",
			Commands: []HelpCommand{
				{Name: "credit-notes list", Summary: "List credit notes.", OptionalFlags: flags("--created-after", "--created-before", "--limit", "--cursor"), Examples: []string{"orb credit-notes list --created-after 2026-04-01T00:00:00Z"}},
				{Name: "credit-notes get", Summary: "Fetch one credit note.", RequiredFlags: []string{"--id"}, Examples: []string{"orb credit-notes get --id cn_123"}},
			},
		},
		"events": {
			Topic: "events",
			Commands: []HelpCommand{
				{Name: "events search", Summary: "Search known events by explicit event IDs.", RequiredFlags: []string{"--id or --ids-file"}, OptionalFlags: flags("--from", "--to"), Examples: []string{"orb events search --id event_id_1 --id event_id_2", "orb events search --ids-file ./event_ids.txt --from 2026-04-01T00:00:00Z"}},
				{Name: "events volume", Summary: "List hourly event ingestion volume.", RequiredFlags: []string{"--from"}, OptionalFlags: flags("--to", "--limit", "--cursor"), Examples: []string{"orb events volume --from 2026-04-01T00:00:00Z --to 2026-04-02T00:00:00Z"}},
				{Name: "events backfills list", Summary: "List event backfills.", OptionalFlags: flags("--limit", "--cursor"), Examples: []string{"orb events backfills list"}},
				{Name: "events backfills get", Summary: "Fetch one event backfill.", RequiredFlags: []string{"--id"}, Examples: []string{"orb events backfills get --id backfill_123"}},
			},
			Notes: []string{
				"events search requires explicit event IDs because Orb's SDK search endpoint is not a general query by customer, event name, or property.",
				"Use events volume when you need ingestion-volume evidence without known event IDs.",
			},
		},
		"alerts": {
			Topic: "alerts",
			Commands: []HelpCommand{
				{Name: "alerts list", Summary: "List alerts.", OptionalFlags: flags("--customer-id", "--external-customer-id", "--subscription-id", "--created-after", "--created-before", "--limit", "--cursor"), Examples: []string{"orb alerts list --customer-id cus_123"}},
				{Name: "alerts get", Summary: "Fetch one alert.", RequiredFlags: []string{"--id"}, Examples: []string{"orb alerts get --id alert_123"}},
			},
		},
		"examples": {
			Topic: "examples",
			Commands: []HelpCommand{
				{Name: "customer investigation", Summary: "Fetch customer, subscriptions, invoices, credits, and alerts.", Examples: []string{"orb customers get --external-id workspace_123", "orb subscriptions list --external-customer-id workspace_123", "orb invoices list --external-customer-id workspace_123", "orb alerts list --external-customer-id workspace_123"}},
				{Name: "usage investigation", Summary: "Fetch subscription usage and event volume for the same timeframe.", Examples: []string{"orb subscriptions usage --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z", "orb events volume --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z"}},
			},
		},
	}
}

func globalHelpFlags() []HelpFlag {
	return []HelpFlag{
		{Name: "--api-key VALUE", Description: "Override ORB_API_KEY for this invocation."},
		{Name: "--base-url URL", Description: "Override ORB_BASE_URL for this invocation."},
		{Name: "--timeout DURATION", Description: "Request timeout, default 60s."},
		{Name: "--limit N", Description: "Page size for list commands, default 20."},
		{Name: "--cursor CURSOR", Description: "Pagination cursor."},
		{Name: "--pretty", Description: "Pretty-print JSON."},
	}
}

func topicLabel(topic string) string {
	labels := map[string]string{
		"customers":     "customers",
		"subscriptions": "subscriptions",
		"plans":         "plans",
		"prices":        "prices",
		"metrics":       "billable metrics",
		"invoices":      "invoices",
		"credit-notes":  "credit notes",
		"events":        "events and event backfills",
		"alerts":        "alerts",
	}
	if label, ok := labels[topic]; ok {
		return label
	}
	return topic
}

func flags(names ...string) []HelpFlag {
	result := make([]HelpFlag, 0, len(names))
	for _, name := range names {
		result = append(result, HelpFlag{Name: name, Description: flagDescription(name)})
	}
	return result
}

func flagDescription(name string) string {
	descriptions := map[string]string{
		"--all":                    "Auto-page through results where supported.",
		"--billable-metric-id":     "Filter usage to a billable metric ID.",
		"--created-after":          "Filter records created at or after an RFC3339 timestamp or YYYY-MM-DD date.",
		"--created-before":         "Filter records created at or before an RFC3339 timestamp or YYYY-MM-DD date.",
		"--currency":               "Currency or custom pricing unit.",
		"--cursor":                 "Pagination cursor.",
		"--entry-status":           "Filter credit ledger entries by status.",
		"--entry-type":             "Filter credit ledger entries by type.",
		"--external-customer-id":   "Orb external customer ID.",
		"--external-id":            "External resource ID.",
		"--external-plan-id":       "External plan ID.",
		"--first-dimension-key":    "First usage dimension key.",
		"--first-dimension-value":  "First usage dimension value.",
		"--from":                   "Start of a timeframe, RFC3339 timestamp or YYYY-MM-DD date.",
		"--granularity":            "Usage granularity such as day.",
		"--group-by":               "Usage grouping key.",
		"--id":                     "Orb resource ID.",
		"--ids-file":               "Path to a newline- or comma-separated event ID file.",
		"--include-all-blocks":     "Include all credit blocks.",
		"--invoice-date-after":     "Filter invoices dated at or after an RFC3339 timestamp or YYYY-MM-DD date.",
		"--invoice-date-before":    "Filter invoices dated at or before an RFC3339 timestamp or YYYY-MM-DD date.",
		"--limit":                  "Page size.",
		"--plan-id":                "Orb plan ID.",
		"--second-dimension-key":   "Second usage dimension key.",
		"--second-dimension-value": "Second usage dimension value.",
		"--start-after":            "Filter schedule phases starting at or after an RFC3339 timestamp or YYYY-MM-DD date.",
		"--start-before":           "Filter schedule phases starting at or before an RFC3339 timestamp or YYYY-MM-DD date.",
		"--status":                 "Resource status filter.",
		"--subscription-id":        "Orb subscription ID.",
		"--to":                     "End of a timeframe, RFC3339 timestamp or YYYY-MM-DD date.",
		"--view-mode":              "Orb view mode, such as cumulative or periodic.",
	}
	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return "See command summary for usage."
}
