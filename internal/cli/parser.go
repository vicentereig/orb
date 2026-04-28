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
	Name                 string
	Resource             string
	Operation            string
	Globals              GlobalOptions
	ID                   string
	IDs                  []string
	IDsFile              string
	ExternalID           string
	CreatedAfter         *time.Time
	CreatedBefore        *time.Time
	From                 *time.Time
	To                   *time.Time
	Currency             string
	ViewMode             string
	IncludeAllBlocks     bool
	EntryType            string
	EntryStatus          string
	CustomerID           string
	ExternalCustomerID   string
	PlanID               string
	ExternalPlanID       string
	Status               string
	Granularity          string
	GroupBy              string
	BillableMetricID     string
	FirstDimensionKey    string
	FirstDimensionValue  string
	SecondDimensionKey   string
	SecondDimensionValue string
	StartAfter           *time.Time
	StartBefore          *time.Time
	SubscriptionID       string
	InvoiceDateAfter     *time.Time
	InvoiceDateBefore    *time.Time
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
	case "subscriptions":
		cmd, err := parseSubscriptions(remaining[1:])
		if err != nil {
			return Command{}, err
		}
		cmd.Name = "subscriptions"
		cmd.Resource = "subscriptions"
		cmd.Globals = globals
		return cmd, nil
	case "plans", "prices", "metrics":
		cmd, err := parseCatalog(remaining[0], remaining[1:])
		if err != nil {
			return Command{}, err
		}
		cmd.Name = remaining[0]
		cmd.Resource = remaining[0]
		cmd.Globals = globals
		return cmd, nil
	case "invoices", "credit-notes":
		cmd, err := parseBilling(remaining[0], remaining[1:])
		if err != nil {
			return Command{}, err
		}
		cmd.Name = remaining[0]
		cmd.Resource = remaining[0]
		cmd.Globals = globals
		return cmd, nil
	case "events":
		cmd, err := parseEvents(remaining[1:])
		if err != nil {
			return Command{}, err
		}
		cmd.Name = "events"
		cmd.Resource = "events"
		cmd.Globals = globals
		return cmd, nil
	case "alerts":
		cmd, err := parseAlerts(remaining[1:])
		if err != nil {
			return Command{}, err
		}
		cmd.Name = "alerts"
		cmd.Resource = "alerts"
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

func parseSubscriptions(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, errors.New("subscriptions requires a subcommand")
	}
	cmd := Command{Operation: args[0]}
	if cmd.Operation != "list" && cmd.Operation != "get" && cmd.Operation != "usage" && cmd.Operation != "costs" && cmd.Operation != "schedule" {
		return Command{}, fmt.Errorf("unknown subscriptions subcommand: %s", cmd.Operation)
	}
	if err := parseSubscriptionFlags(&cmd, args[1:]); err != nil {
		return Command{}, err
	}
	switch cmd.Operation {
	case "list":
		return cmd, nil
	case "get", "schedule":
		if cmd.ID == "" {
			return Command{}, fmt.Errorf("subscriptions %s requires --id", cmd.Operation)
		}
	case "usage", "costs":
		if cmd.ID == "" {
			return Command{}, fmt.Errorf("subscriptions %s requires --id", cmd.Operation)
		}
		if cmd.From == nil {
			return Command{}, fmt.Errorf("subscriptions %s requires --from", cmd.Operation)
		}
		if cmd.To == nil {
			return Command{}, fmt.Errorf("subscriptions %s requires --to", cmd.Operation)
		}
		if !cmd.From.Before(*cmd.To) {
			return Command{}, errors.New("--from must be before --to")
		}
	}
	return cmd, nil
}

func parseSubscriptionFlags(cmd *Command, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--id", "--customer-id", "--external-customer-id", "--plan-id", "--external-plan-id", "--status", "--created-after", "--created-before", "--from", "--to", "--currency", "--view-mode", "--granularity", "--group-by", "--billable-metric-id", "--first-dimension-key", "--first-dimension-value", "--second-dimension-key", "--second-dimension-value", "--start-after", "--start-before":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			if err := setSubscriptionFlag(cmd, arg, args[i+1]); err != nil {
				return err
			}
			i++
		default:
			if name, value, ok := strings.Cut(arg, "="); ok {
				if err := setSubscriptionFlag(cmd, name, value); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown subscriptions option: %s", arg)
		}
	}
	return nil
}

func setSubscriptionFlag(cmd *Command, name, value string) error {
	switch name {
	case "--id":
		cmd.ID = value
	case "--customer-id":
		cmd.CustomerID = value
	case "--external-customer-id":
		cmd.ExternalCustomerID = value
	case "--plan-id":
		cmd.PlanID = value
	case "--external-plan-id":
		cmd.ExternalPlanID = value
	case "--status":
		cmd.Status = value
	case "--currency":
		cmd.Currency = value
	case "--view-mode":
		cmd.ViewMode = value
	case "--granularity":
		cmd.Granularity = value
	case "--group-by":
		cmd.GroupBy = value
	case "--billable-metric-id":
		cmd.BillableMetricID = value
	case "--first-dimension-key":
		cmd.FirstDimensionKey = value
	case "--first-dimension-value":
		cmd.FirstDimensionValue = value
	case "--second-dimension-key":
		cmd.SecondDimensionKey = value
	case "--second-dimension-value":
		cmd.SecondDimensionValue = value
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
	case "--start-after":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --start-after: %w", err)
		}
		cmd.StartAfter = &t
	case "--start-before":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --start-before: %w", err)
		}
		cmd.StartBefore = &t
	default:
		return fmt.Errorf("unknown subscriptions option: %s", name)
	}
	return nil
}

func parseCatalog(resource string, args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, fmt.Errorf("%s requires a subcommand", resource)
	}
	cmd := Command{Operation: args[0]}
	if cmd.Operation != "list" && cmd.Operation != "get" {
		return Command{}, fmt.Errorf("unknown %s subcommand: %s", resource, cmd.Operation)
	}
	if err := parseCatalogFlags(resource, &cmd, args[1:]); err != nil {
		return Command{}, err
	}
	if cmd.Operation == "get" {
		switch resource {
		case "metrics":
			if cmd.ExternalID != "" {
				return Command{}, errors.New("metrics get does not support --external-id")
			}
			if cmd.ID == "" {
				return Command{}, errors.New("metrics get requires --id")
			}
		default:
			if err := validateOneIdentity(cmd); err != nil {
				return Command{}, fmt.Errorf("%s get: %w", resource, err)
			}
		}
	}
	return cmd, nil
}

func parseCatalogFlags(resource string, cmd *Command, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--id", "--external-id", "--created-after", "--created-before":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			if err := setCatalogFlag(resource, cmd, arg, args[i+1]); err != nil {
				return err
			}
			i++
		case "--status":
			if resource != "plans" {
				return fmt.Errorf("%s does not support --status", resource)
			}
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			cmd.Status = args[i+1]
			i++
		default:
			if name, value, ok := strings.Cut(arg, "="); ok {
				if err := setCatalogFlag(resource, cmd, name, value); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown %s option: %s", resource, arg)
		}
	}
	return nil
}

func setCatalogFlag(resource string, cmd *Command, name, value string) error {
	switch name {
	case "--id":
		cmd.ID = value
	case "--external-id":
		if resource == "metrics" {
			return errors.New("metrics does not support --external-id")
		}
		cmd.ExternalID = value
	case "--created-after":
		if resource == "prices" {
			return errors.New("prices does not support --created-after")
		}
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --created-after: %w", err)
		}
		cmd.CreatedAfter = &t
	case "--created-before":
		if resource == "prices" {
			return errors.New("prices does not support --created-before")
		}
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --created-before: %w", err)
		}
		cmd.CreatedBefore = &t
	case "--status":
		if resource != "plans" {
			return fmt.Errorf("%s does not support --status", resource)
		}
		cmd.Status = value
	default:
		return fmt.Errorf("unknown %s option: %s", resource, name)
	}
	return nil
}

func parseBilling(resource string, args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, fmt.Errorf("%s requires a subcommand", resource)
	}
	cmd := Command{Operation: args[0]}
	valid := map[string]bool{"list": true, "get": true}
	if resource == "invoices" {
		valid["summary"] = true
		valid["upcoming"] = true
	}
	if !valid[cmd.Operation] {
		return Command{}, fmt.Errorf("unknown %s subcommand: %s", resource, cmd.Operation)
	}
	if err := parseBillingFlags(resource, &cmd, args[1:]); err != nil {
		return Command{}, err
	}
	switch cmd.Operation {
	case "get":
		if cmd.ID == "" {
			return Command{}, fmt.Errorf("%s get requires --id", resource)
		}
	case "upcoming":
		if cmd.SubscriptionID == "" {
			return Command{}, errors.New("invoices upcoming requires --subscription-id")
		}
	}
	return cmd, nil
}

func parseBillingFlags(resource string, cmd *Command, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--id", "--created-after", "--created-before", "--customer-id", "--external-customer-id", "--subscription-id", "--status", "--invoice-date-after", "--invoice-date-before":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			if err := setBillingFlag(resource, cmd, arg, args[i+1]); err != nil {
				return err
			}
			i++
		default:
			if name, value, ok := strings.Cut(arg, "="); ok {
				if err := setBillingFlag(resource, cmd, name, value); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown %s option: %s", resource, arg)
		}
	}
	return nil
}

func setBillingFlag(resource string, cmd *Command, name, value string) error {
	if resource == "credit-notes" {
		switch name {
		case "--customer-id", "--external-customer-id", "--subscription-id", "--status", "--invoice-date-after", "--invoice-date-before":
			return fmt.Errorf("credit-notes does not support %s", name)
		}
	}
	switch name {
	case "--id":
		cmd.ID = value
	case "--customer-id":
		cmd.CustomerID = value
	case "--external-customer-id":
		cmd.ExternalCustomerID = value
	case "--subscription-id":
		cmd.SubscriptionID = value
	case "--status":
		cmd.Status = value
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
	case "--invoice-date-after":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --invoice-date-after: %w", err)
		}
		cmd.InvoiceDateAfter = &t
	case "--invoice-date-before":
		t, err := parseTime(value)
		if err != nil {
			return fmt.Errorf("invalid --invoice-date-before: %w", err)
		}
		cmd.InvoiceDateBefore = &t
	default:
		return fmt.Errorf("unknown %s option: %s", resource, name)
	}
	return nil
}

func parseEvents(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, errors.New("events requires a subcommand")
	}
	if args[0] == "backfills" {
		return parseEventBackfills(args[1:])
	}
	cmd := Command{Operation: args[0]}
	if cmd.Operation != "search" && cmd.Operation != "volume" {
		return Command{}, fmt.Errorf("unknown events subcommand: %s", cmd.Operation)
	}
	if err := parseEventFlags(&cmd, args[1:]); err != nil {
		return Command{}, err
	}
	switch cmd.Operation {
	case "search":
		if len(cmd.IDs) == 0 && cmd.IDsFile == "" {
			return Command{}, errors.New("events search requires --id or --ids-file")
		}
	case "volume":
		if cmd.From == nil {
			return Command{}, errors.New("events volume requires --from")
		}
	}
	if cmd.From != nil && cmd.To != nil && !cmd.From.Before(*cmd.To) {
		return Command{}, errors.New("--from must be before --to")
	}
	return cmd, nil
}

func parseEventBackfills(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, errors.New("events backfills requires a subcommand")
	}
	cmd := Command{Operation: "backfills-" + args[0]}
	if cmd.Operation != "backfills-list" && cmd.Operation != "backfills-get" {
		return Command{}, fmt.Errorf("unknown events backfills subcommand: %s", args[0])
	}
	if err := parseEventFlags(&cmd, args[1:]); err != nil {
		return Command{}, err
	}
	if cmd.Operation == "backfills-get" && cmd.ID == "" {
		return Command{}, errors.New("events backfills get requires --id")
	}
	return cmd, nil
}

func parseEventFlags(cmd *Command, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--id", "--ids-file", "--from", "--to":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			if err := setEventFlag(cmd, arg, args[i+1]); err != nil {
				return err
			}
			i++
		default:
			if name, value, ok := strings.Cut(arg, "="); ok {
				if err := setEventFlag(cmd, name, value); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown events option: %s", arg)
		}
	}
	return nil
}

func setEventFlag(cmd *Command, name, value string) error {
	switch name {
	case "--id":
		if cmd.Operation == "search" {
			cmd.IDs = append(cmd.IDs, value)
		} else {
			cmd.ID = value
		}
	case "--ids-file":
		cmd.IDsFile = value
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
	default:
		return fmt.Errorf("unknown events option: %s", name)
	}
	return nil
}

func parseAlerts(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, errors.New("alerts requires a subcommand")
	}
	cmd := Command{Operation: args[0]}
	if cmd.Operation != "list" && cmd.Operation != "get" {
		return Command{}, fmt.Errorf("unknown alerts subcommand: %s", cmd.Operation)
	}
	if err := parseAlertFlags(&cmd, args[1:]); err != nil {
		return Command{}, err
	}
	if cmd.Operation == "get" && cmd.ID == "" {
		return Command{}, errors.New("alerts get requires --id")
	}
	return cmd, nil
}

func parseAlertFlags(cmd *Command, args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--id", "--customer-id", "--external-customer-id", "--subscription-id", "--created-after", "--created-before":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			if err := setAlertFlag(cmd, arg, args[i+1]); err != nil {
				return err
			}
			i++
		default:
			if name, value, ok := strings.Cut(arg, "="); ok {
				if err := setAlertFlag(cmd, name, value); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown alerts option: %s", arg)
		}
	}
	return nil
}

func setAlertFlag(cmd *Command, name, value string) error {
	switch name {
	case "--id":
		cmd.ID = value
	case "--customer-id":
		cmd.CustomerID = value
	case "--external-customer-id":
		cmd.ExternalCustomerID = value
	case "--subscription-id":
		cmd.SubscriptionID = value
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
	default:
		return fmt.Errorf("unknown alerts option: %s", name)
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
