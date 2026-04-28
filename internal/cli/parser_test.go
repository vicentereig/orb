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
