package cli

import (
	"testing"
	"time"
)

func TestParseVersion(t *testing.T) {
	cmd, err := Parse([]string{"version"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if cmd.Name != "version" || cmd.Operation != "version" {
		t.Fatalf("unexpected command: %+v", cmd)
	}
}

func TestParsePingWithGlobalsAnywhere(t *testing.T) {
	cmd, err := Parse([]string{"--pretty", "ping", "--timeout=10s", "--limit", "50", "--base-url", "http://localhost:8080"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if cmd.Name != "ping" {
		t.Fatalf("expected ping, got %+v", cmd)
	}
	if !cmd.Globals.Pretty {
		t.Fatalf("expected pretty global")
	}
	if cmd.Globals.Timeout != 10*time.Second {
		t.Fatalf("unexpected timeout: %s", cmd.Globals.Timeout)
	}
	if cmd.Globals.Limit != 50 {
		t.Fatalf("unexpected limit: %d", cmd.Globals.Limit)
	}
	if cmd.Globals.BaseURL != "http://localhost:8080" {
		t.Fatalf("unexpected base URL: %s", cmd.Globals.BaseURL)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "missing command", args: nil},
		{name: "unknown command", args: []string{"customers"}},
		{name: "missing timeout value", args: []string{"--timeout"}},
		{name: "invalid limit", args: []string{"ping", "--limit", "zero"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestParseCustomerCommands(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOp     string
		wantID     string
		wantExtID  string
		wantCursor string
	}{
		{
			name:       "customers list",
			args:       []string{"customers", "list", "--created-after", "2026-04-01T00:00:00Z", "--created-before", "2026-05-01T00:00:00Z", "--cursor", "next"},
			wantOp:     "list",
			wantCursor: "next",
		},
		{
			name:   "customers get by id",
			args:   []string{"customers", "get", "--id", "cus_123"},
			wantOp: "get",
			wantID: "cus_123",
		},
		{
			name:      "customers get by external id",
			args:      []string{"customers", "get", "--external-id", "workspace_123"},
			wantOp:    "get",
			wantExtID: "workspace_123",
		},
		{
			name:   "customer costs",
			args:   []string{"customers", "costs", "--id", "cus_123", "--from", "2026-04-01T00:00:00Z", "--to", "2026-05-01T00:00:00Z", "--currency", "USD", "--view-mode", "periodic"},
			wantOp: "costs",
			wantID: "cus_123",
		},
		{
			name:      "customer credits by external id",
			args:      []string{"customers", "credits", "--external-id", "workspace_123", "--include-all-blocks"},
			wantOp:    "credits",
			wantExtID: "workspace_123",
		},
		{
			name:   "customer credit ledger",
			args:   []string{"customers", "credit-ledger", "--id", "cus_123", "--entry-type", "increment", "--entry-status", "committed"},
			wantOp: "credit-ledger",
			wantID: "cus_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if cmd.Resource != "customers" || cmd.Operation != tt.wantOp {
				t.Fatalf("unexpected command: %+v", cmd)
			}
			if cmd.ID != tt.wantID || cmd.ExternalID != tt.wantExtID {
				t.Fatalf("unexpected identity: %+v", cmd)
			}
			if cmd.Globals.Cursor != tt.wantCursor {
				t.Fatalf("unexpected cursor: %q", cmd.Globals.Cursor)
			}
		})
	}
}

func TestParseCustomerCommandValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "customers requires subcommand", args: []string{"customers"}},
		{name: "customers get requires id", args: []string{"customers", "get"}},
		{name: "customers get rejects both ids", args: []string{"customers", "get", "--id", "cus_123", "--external-id", "workspace_123"}},
		{name: "costs requires from", args: []string{"customers", "costs", "--id", "cus_123", "--to", "2026-05-01T00:00:00Z"}},
		{name: "costs requires to", args: []string{"customers", "costs", "--id", "cus_123", "--from", "2026-04-01T00:00:00Z"}},
		{name: "from must be before to", args: []string{"customers", "costs", "--id", "cus_123", "--from", "2026-05-01T00:00:00Z", "--to", "2026-04-01T00:00:00Z"}},
		{name: "credits requires id", args: []string{"customers", "credits"}},
		{name: "credit-ledger rejects unknown flag", args: []string{"customers", "credit-ledger", "--id", "cus_123", "--wat"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestParseSubscriptionCommands(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOp     string
		wantID     string
		wantStatus string
	}{
		{
			name:       "subscriptions list",
			args:       []string{"subscriptions", "list", "--customer-id", "cus_123", "--external-plan-id", "pro", "--status", "active"},
			wantOp:     "list",
			wantStatus: "active",
		},
		{
			name:   "subscriptions get",
			args:   []string{"subscriptions", "get", "--id", "sub_123"},
			wantOp: "get",
			wantID: "sub_123",
		},
		{
			name:   "subscriptions usage",
			args:   []string{"subscriptions", "usage", "--id", "sub_123", "--from", "2026-04-01T00:00:00Z", "--to", "2026-05-01T00:00:00Z", "--granularity", "day", "--group-by", "region", "--billable-metric-id", "bm_123"},
			wantOp: "usage",
			wantID: "sub_123",
		},
		{
			name:   "subscriptions costs",
			args:   []string{"subscriptions", "costs", "--id", "sub_123", "--from", "2026-04-01T00:00:00Z", "--to", "2026-05-01T00:00:00Z", "--currency", "USD", "--view-mode", "cumulative"},
			wantOp: "costs",
			wantID: "sub_123",
		},
		{
			name:   "subscriptions schedule",
			args:   []string{"subscriptions", "schedule", "--id", "sub_123", "--start-after", "2026-04-01T00:00:00Z"},
			wantOp: "schedule",
			wantID: "sub_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if cmd.Resource != "subscriptions" || cmd.Operation != tt.wantOp {
				t.Fatalf("unexpected command: %+v", cmd)
			}
			if cmd.ID != tt.wantID {
				t.Fatalf("unexpected id: %+v", cmd)
			}
			if cmd.Status != tt.wantStatus {
				t.Fatalf("unexpected status: %+v", cmd)
			}
		})
	}
}

func TestParseSubscriptionCommandValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "subscriptions requires subcommand", args: []string{"subscriptions"}},
		{name: "get requires id", args: []string{"subscriptions", "get"}},
		{name: "usage requires id", args: []string{"subscriptions", "usage", "--from", "2026-04-01T00:00:00Z", "--to", "2026-05-01T00:00:00Z"}},
		{name: "usage requires from", args: []string{"subscriptions", "usage", "--id", "sub_123", "--to", "2026-05-01T00:00:00Z"}},
		{name: "usage requires to", args: []string{"subscriptions", "usage", "--id", "sub_123", "--from", "2026-04-01T00:00:00Z"}},
		{name: "costs rejects bad time range", args: []string{"subscriptions", "costs", "--id", "sub_123", "--from", "2026-05-01T00:00:00Z", "--to", "2026-04-01T00:00:00Z"}},
		{name: "schedule requires id", args: []string{"subscriptions", "schedule"}},
		{name: "unknown flag", args: []string{"subscriptions", "list", "--wat"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestParseCatalogCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		resource string
		op       string
		id       string
		extID    string
		status   string
	}{
		{name: "plans list", args: []string{"plans", "list", "--status", "active"}, resource: "plans", op: "list", status: "active"},
		{name: "plans get", args: []string{"plans", "get", "--id", "plan_123"}, resource: "plans", op: "get", id: "plan_123"},
		{name: "plans get external", args: []string{"plans", "get", "--external-id", "pro"}, resource: "plans", op: "get", extID: "pro"},
		{name: "prices list", args: []string{"prices", "list", "--limit", "50"}, resource: "prices", op: "list"},
		{name: "prices get external", args: []string{"prices", "get", "--external-id", "api_calls"}, resource: "prices", op: "get", extID: "api_calls"},
		{name: "metrics list", args: []string{"metrics", "list", "--created-after", "2026-04-01T00:00:00Z"}, resource: "metrics", op: "list"},
		{name: "metrics get", args: []string{"metrics", "get", "--id", "metric_123"}, resource: "metrics", op: "get", id: "metric_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if cmd.Resource != tt.resource || cmd.Operation != tt.op {
				t.Fatalf("unexpected command: %+v", cmd)
			}
			if cmd.ID != tt.id || cmd.ExternalID != tt.extID || cmd.Status != tt.status {
				t.Fatalf("unexpected command fields: %+v", cmd)
			}
		})
	}
}

func TestParseCatalogCommandValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "plans requires subcommand", args: []string{"plans"}},
		{name: "plans get requires id", args: []string{"plans", "get"}},
		{name: "metrics get rejects external id", args: []string{"metrics", "get", "--external-id", "metric"}},
		{name: "prices unknown flag", args: []string{"prices", "list", "--status", "active"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestParseBillingCommands(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		resource       string
		op             string
		id             string
		subscriptionID string
	}{
		{name: "invoices list", args: []string{"invoices", "list", "--customer-id", "cus_123", "--subscription-id", "sub_123", "--status", "issued", "--invoice-date-after", "2026-04-01T00:00:00Z"}, resource: "invoices", op: "list", subscriptionID: "sub_123"},
		{name: "invoices get", args: []string{"invoices", "get", "--id", "inv_123"}, resource: "invoices", op: "get", id: "inv_123"},
		{name: "invoices summary", args: []string{"invoices", "summary", "--external-customer-id", "workspace_123"}, resource: "invoices", op: "summary"},
		{name: "invoices upcoming", args: []string{"invoices", "upcoming", "--subscription-id", "sub_123"}, resource: "invoices", op: "upcoming", subscriptionID: "sub_123"},
		{name: "credit notes list", args: []string{"credit-notes", "list", "--created-after", "2026-04-01T00:00:00Z"}, resource: "credit-notes", op: "list"},
		{name: "credit notes get", args: []string{"credit-notes", "get", "--id", "cn_123"}, resource: "credit-notes", op: "get", id: "cn_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if cmd.Resource != tt.resource || cmd.Operation != tt.op || cmd.ID != tt.id || cmd.SubscriptionID != tt.subscriptionID {
				t.Fatalf("unexpected command: %+v", cmd)
			}
		})
	}
}

func TestParseBillingCommandValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "invoices requires subcommand", args: []string{"invoices"}},
		{name: "invoices get requires id", args: []string{"invoices", "get"}},
		{name: "invoices upcoming requires subscription id", args: []string{"invoices", "upcoming"}},
		{name: "credit notes get requires id", args: []string{"credit-notes", "get"}},
		{name: "credit notes rejects customer filter", args: []string{"credit-notes", "list", "--customer-id", "cus_123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}

func TestParseEventAndAlertCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		resource string
		op       string
		id       string
		ids      int
	}{
		{name: "events search", args: []string{"events", "search", "--id", "evt_1", "--id", "evt_2", "--from", "2026-04-01T00:00:00Z"}, resource: "events", op: "search", ids: 2},
		{name: "events volume", args: []string{"events", "volume", "--from", "2026-04-01T00:00:00Z", "--to", "2026-04-02T00:00:00Z"}, resource: "events", op: "volume"},
		{name: "events backfills list", args: []string{"events", "backfills", "list"}, resource: "events", op: "backfills-list"},
		{name: "events backfills get", args: []string{"events", "backfills", "get", "--id", "bf_123"}, resource: "events", op: "backfills-get", id: "bf_123"},
		{name: "alerts list", args: []string{"alerts", "list", "--customer-id", "cus_123"}, resource: "alerts", op: "list"},
		{name: "alerts get", args: []string{"alerts", "get", "--id", "alert_123"}, resource: "alerts", op: "get", id: "alert_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if cmd.Resource != tt.resource || cmd.Operation != tt.op || cmd.ID != tt.id || len(cmd.IDs) != tt.ids {
				t.Fatalf("unexpected command: %+v", cmd)
			}
		})
	}
}

func TestParseEventAndAlertValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "events search requires id or ids file", args: []string{"events", "search"}},
		{name: "events volume requires from", args: []string{"events", "volume"}},
		{name: "backfills get requires id", args: []string{"events", "backfills", "get"}},
		{name: "alerts get requires id", args: []string{"alerts", "get"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Fatalf("expected error")
			}
		})
	}
}
