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
