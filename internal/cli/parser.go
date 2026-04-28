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
	Name      string
	Resource  string
	Operation string
	Globals   GlobalOptions
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
	default:
		return Command{}, fmt.Errorf("unknown command: %s", remaining[0])
	}
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
