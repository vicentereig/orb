package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GlobalOptions struct {
	APIKey   string
	BaseURL  string
	Timeout  time.Duration
	Limit    int
	Cursor   string
	All      bool
	Pretty   bool
	Raw      bool
	Verbose  bool
	MaxPages int
}

type Command struct {
	Name             string
	Resource         string
	Operation        string
	Globals          GlobalOptions
	ID               string
	ExternalID       string
	CreatedAfter     *time.Time
	CreatedBefore    *time.Time
	From             *time.Time
	To               *time.Time
	Currency         string
	ViewMode         string
	IncludeAllBlocks bool
	EntryType        string
	EntryStatus      string
}

func Parse(args []string) (Command, error) {
	globals := GlobalOptions{
		Timeout: 60 * time.Second,
		Limit:   20,
	}

	remaining := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--pretty":
			globals.Pretty = true
		case arg == "--raw":
			globals.Raw = true
		case arg == "--verbose":
			globals.Verbose = true
		case arg == "--all":
			globals.All = true
		case arg == "--api-key" || arg == "--base-url" || arg == "--timeout" || arg == "--limit" || arg == "--cursor" || arg == "--max-pages":
			if i+1 >= len(args) {
				return Command{}, fmt.Errorf("%s requires a value", arg)
			}
			if err := setGlobal(&globals, arg, args[i+1]); err != nil {
				return Command{}, err
			}
			i++
		case strings.HasPrefix(arg, "--api-key=") || strings.HasPrefix(arg, "--base-url=") ||
			strings.HasPrefix(arg, "--timeout=") || strings.HasPrefix(arg, "--limit=") ||
			strings.HasPrefix(arg, "--cursor=") || strings.HasPrefix(arg, "--max-pages="):
			name, value, _ := strings.Cut(arg, "=")
			if err := setGlobal(&globals, name, value); err != nil {
				return Command{}, err
			}
		default:
			remaining = append(remaining, arg)
		}
	}

	if len(remaining) == 0 {
		return Command{}, errors.New("command required")
	}

	switch remaining[0] {
	case "version":
		return Command{Name: "version", Operation: "version", Globals: globals}, nil
	case "ping":
		return Command{Name: "ping", Operation: "ping", Globals: globals}, nil
	case "customers":
		cmd, err := parseCustomers(remaining[1:])
		if err != nil {
			return Command{}, err
		}
		cmd.Name = "customers"
		cmd.Resource = "customers"
		cmd.Globals = globals
		return cmd, nil
	default:
		return Command{}, fmt.Errorf("unknown command: %s", remaining[0])
	}
}

func parseCustomers(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, errors.New("customers requires a subcommand")
	}
	cmd := Command{Operation: args[0]}
	if cmd.Operation != "list" && cmd.Operation != "get" && cmd.Operation != "costs" && cmd.Operation != "credits" && cmd.Operation != "credit-ledger" {
		return Command{}, fmt.Errorf("unknown customers subcommand: %s", cmd.Operation)
	}
	if err := parseCustomerFlags(&cmd, args[1:]); err != nil {
		return Command{}, err
	}
	switch cmd.Operation {
	case "list":
		return cmd, nil
	case "get", "credits", "credit-ledger":
		if err := validateOneIdentity(cmd); err != nil {
			return Command{}, err
		}
	case "costs":
		if err := validateOneIdentity(cmd); err != nil {
			return Command{}, err
		}
		if cmd.From == nil {
			return Command{}, errors.New("customers costs requires --from")
		}
		if cmd.To == nil {
			return Command{}, errors.New("customers costs requires --to")
		}
		if !cmd.From.Before(*cmd.To) {
			return Command{}, errors.New("--from must be before --to")
		}
	}
	return cmd, nil
}

func parseCustomerFlags(cmd *Command, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--include-all-blocks":
			cmd.IncludeAllBlocks = true
		case "--id", "--external-id", "--created-after", "--created-before", "--from", "--to", "--currency", "--view-mode", "--entry-type", "--entry-status":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			if err := setCustomerFlag(cmd, arg, args[i+1]); err != nil {
				return err
			}
			i++
		default:
			if name, value, ok := strings.Cut(arg, "="); ok {
				if err := setCustomerFlag(cmd, name, value); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown customers option: %s", arg)
		}
	}
	return nil
}

func setCustomerFlag(cmd *Command, name, value string) error {
	switch name {
	case "--id":
		cmd.ID = value
	case "--external-id":
		cmd.ExternalID = value
	case "--created-after":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --created-after: %w", err)
		}
		cmd.CreatedAfter = &t
	case "--created-before":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --created-before: %w", err)
		}
		cmd.CreatedBefore = &t
	case "--from":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --from: %w", err)
		}
		cmd.From = &t
	case "--to":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --to: %w", err)
		}
		cmd.To = &t
	case "--currency":
		cmd.Currency = value
	case "--view-mode":
		cmd.ViewMode = value
	case "--entry-type":
		cmd.EntryType = value
	case "--entry-status":
		cmd.EntryStatus = value
	default:
		return fmt.Errorf("unknown customers option: %s", name)
	}
	return nil
}

func parseTime(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", value)
}

func validateOneIdentity(cmd Command) error {
	if cmd.ID == "" && cmd.ExternalID == "" {
		return fmt.Errorf("customers %s requires --id or --external-id", cmd.Operation)
	}
	if cmd.ID != "" && cmd.ExternalID != "" {
		return fmt.Errorf("customers %s accepts only one of --id or --external-id", cmd.Operation)
	}
	return nil
}

func setGlobal(globals *GlobalOptions, name, value string) error {
	switch name {
	case "--api-key":
		globals.APIKey = value
	case "--base-url":
		globals.BaseURL = value
	case "--timeout":
		timeout, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid --timeout: %w", err)
		}
		globals.Timeout = timeout
	case "--limit":
		limit, err := strconv.Atoi(value)
		if err != nil || limit < 1 {
			return fmt.Errorf("invalid --limit: %s", value)
		}
		globals.Limit = limit
	case "--cursor":
		globals.Cursor = value
	case "--max-pages":
		maxPages, err := strconv.Atoi(value)
		if err != nil || maxPages < 1 {
			return fmt.Errorf("invalid --max-pages: %s", value)
		}
		globals.MaxPages = maxPages
	default:
		return fmt.Errorf("unknown global option: %s", name)
	}
	return nil
}
