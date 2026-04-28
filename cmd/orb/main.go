package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/vicentereig/orb/internal/app"
	"github.com/vicentereig/orb/internal/cli"
	"github.com/vicentereig/orb/internal/orbclient"
	"github.com/vicentereig/orb/internal/output"
)

var version = "dev"

const usage = `orb - Orb billing forensics from the terminal

Usage:
  orb <command> [options]

Commands:
  version                    Print CLI version information
  ping                       Check Orb API connectivity and credentials
  customers list|get|costs|credits|credit-ledger
  subscriptions list|get|usage|costs|schedule
  plans list|get
  prices list|get
  metrics list|get
  invoices list|get|summary|upcoming
  credit-notes list|get
  events search|volume
  events backfills list|get
  alerts list|get

Global options:
  --api-key VALUE     Override ORB_API_KEY for this invocation
  --base-url URL      Override ORB_BASE_URL for this invocation
  --timeout DURATION  Request timeout (default: 60s)
  --limit N           Page size for list commands (default: 20)
  --cursor CURSOR     Pagination cursor
  --pretty            Pretty-print JSON
`

func main() {
	cmd, err := cli.Parse(os.Args[1:])
	if err != nil {
		if len(os.Args) == 1 {
			fmt.Fprint(os.Stderr, usage)
		}
		fmt.Fprintln(os.Stderr, output.Error(err, "usage_error", nil))
		os.Exit(1)
	}

	client := orbclient.New(cmd.Globals.APIKey, cmd.Globals.BaseURL)
	application := app.New(client, version)

	ctx, cancel := context.WithTimeout(context.Background(), cmd.Globals.Timeout)
	defer cancel()

	var result string
	switch cmd.Name {
	case "version":
		result = application.Version()
	case "ping":
		result = application.Ping(ctx)
	case "customers":
		result = executeCustomers(ctx, application, cmd)
	case "subscriptions":
		result = executeSubscriptions(ctx, application, cmd)
	case "plans", "prices", "metrics":
		result = executeCatalog(ctx, application, cmd)
	case "invoices", "credit-notes":
		result = executeBilling(ctx, application, cmd)
	case "events":
		result = executeEvents(ctx, application, cmd)
	case "alerts":
		result = executeAlerts(ctx, application, cmd)
	default:
		result = output.Error(fmt.Errorf("unknown command: %s", cmd.Name), "usage_error", nil)
	}

	if cmd.Globals.Pretty {
		if pretty, err := output.Pretty(result); err == nil {
			result = pretty
		}
	}

	fmt.Println(result)
}

func executeCustomers(ctx context.Context, application *app.App, cmd cli.Command) string {
	identity := app.CustomerIdentity{ID: cmd.ID, ExternalID: cmd.ExternalID}
	switch cmd.Operation {
	case "list":
		return application.ListCustomers(ctx, app.CustomerListParams{
			Limit:         cmd.Globals.Limit,
			Cursor:        cmd.Globals.Cursor,
			CreatedAfter:  cmd.CreatedAfter,
			CreatedBefore: cmd.CreatedBefore,
		})
	case "get":
		return application.GetCustomer(ctx, identity)
	case "costs":
		return application.ListCustomerCosts(ctx, identity, app.TimeframeParams{
			From:     cmd.From,
			To:       cmd.To,
			Currency: cmd.Currency,
			ViewMode: cmd.ViewMode,
		})
	case "credits":
		return application.ListCustomerCredits(ctx, identity, app.CustomerCreditsParams{
			Limit:            cmd.Globals.Limit,
			Cursor:           cmd.Globals.Cursor,
			Currency:         cmd.Currency,
			IncludeAllBlocks: cmd.IncludeAllBlocks,
		})
	case "credit-ledger":
		return application.ListCustomerCreditLedger(ctx, identity, app.CustomerCreditLedgerParams{
			Limit:       cmd.Globals.Limit,
			Cursor:      cmd.Globals.Cursor,
			Currency:    cmd.Currency,
			EntryType:   cmd.EntryType,
			EntryStatus: cmd.EntryStatus,
		})
	default:
		return output.Error(fmt.Errorf("unknown customers subcommand: %s", cmd.Operation), "usage_error", nil)
	}
}

func executeSubscriptions(ctx context.Context, application *app.App, cmd cli.Command) string {
	switch cmd.Operation {
	case "list":
		return application.ListSubscriptions(ctx, app.SubscriptionListParams{
			Limit:              cmd.Globals.Limit,
			Cursor:             cmd.Globals.Cursor,
			CreatedAfter:       cmd.CreatedAfter,
			CreatedBefore:      cmd.CreatedBefore,
			CustomerID:         cmd.CustomerID,
			ExternalCustomerID: cmd.ExternalCustomerID,
			PlanID:             cmd.PlanID,
			ExternalPlanID:     cmd.ExternalPlanID,
			Status:             cmd.Status,
		})
	case "get":
		return application.GetSubscription(ctx, cmd.ID)
	case "usage":
		return application.ListSubscriptionUsage(ctx, cmd.ID, app.SubscriptionUsageParams{
			From:                 cmd.From,
			To:                   cmd.To,
			BillableMetricID:     cmd.BillableMetricID,
			FirstDimensionKey:    cmd.FirstDimensionKey,
			FirstDimensionValue:  cmd.FirstDimensionValue,
			Granularity:          cmd.Granularity,
			GroupBy:              cmd.GroupBy,
			SecondDimensionKey:   cmd.SecondDimensionKey,
			SecondDimensionValue: cmd.SecondDimensionValue,
			ViewMode:             cmd.ViewMode,
		})
	case "costs":
		return application.ListSubscriptionCosts(ctx, cmd.ID, app.TimeframeParams{
			From:     cmd.From,
			To:       cmd.To,
			Currency: cmd.Currency,
			ViewMode: cmd.ViewMode,
		})
	case "schedule":
		return application.ListSubscriptionSchedule(ctx, cmd.ID, app.SubscriptionScheduleParams{
			Limit:       cmd.Globals.Limit,
			Cursor:      cmd.Globals.Cursor,
			StartAfter:  cmd.StartAfter,
			StartBefore: cmd.StartBefore,
		})
	default:
		return output.Error(fmt.Errorf("unknown subscriptions subcommand: %s", cmd.Operation), "usage_error", nil)
	}
}

func executeCatalog(ctx context.Context, application *app.App, cmd cli.Command) string {
	switch cmd.Resource {
	case "plans":
		switch cmd.Operation {
		case "list":
			return application.ListPlans(ctx, app.CatalogListParams{
				Limit:         cmd.Globals.Limit,
				Cursor:        cmd.Globals.Cursor,
				CreatedAfter:  cmd.CreatedAfter,
				CreatedBefore: cmd.CreatedBefore,
				Status:        cmd.Status,
			})
		case "get":
			return application.GetPlan(ctx, app.ResourceIdentity{ID: cmd.ID, ExternalID: cmd.ExternalID})
		}
	case "prices":
		switch cmd.Operation {
		case "list":
			return application.ListPrices(ctx, app.PageParams{Limit: cmd.Globals.Limit, Cursor: cmd.Globals.Cursor})
		case "get":
			return application.GetPrice(ctx, app.ResourceIdentity{ID: cmd.ID, ExternalID: cmd.ExternalID})
		}
	case "metrics":
		switch cmd.Operation {
		case "list":
			return application.ListMetrics(ctx, app.CatalogListParams{
				Limit:         cmd.Globals.Limit,
				Cursor:        cmd.Globals.Cursor,
				CreatedAfter:  cmd.CreatedAfter,
				CreatedBefore: cmd.CreatedBefore,
			})
		case "get":
			return application.GetMetric(ctx, cmd.ID)
		}
	}
	return output.Error(fmt.Errorf("unknown %s subcommand: %s", cmd.Resource, cmd.Operation), "usage_error", nil)
}

func executeBilling(ctx context.Context, application *app.App, cmd cli.Command) string {
	switch cmd.Resource {
	case "invoices":
		params := app.InvoiceListParams{
			Limit:              cmd.Globals.Limit,
			Cursor:             cmd.Globals.Cursor,
			CustomerID:         cmd.CustomerID,
			ExternalCustomerID: cmd.ExternalCustomerID,
			SubscriptionID:     cmd.SubscriptionID,
			Status:             cmd.Status,
			InvoiceDateAfter:   cmd.InvoiceDateAfter,
			InvoiceDateBefore:  cmd.InvoiceDateBefore,
		}
		switch cmd.Operation {
		case "list":
			return application.ListInvoices(ctx, params)
		case "get":
			return application.GetInvoice(ctx, cmd.ID)
		case "summary":
			return application.ListInvoiceSummary(ctx, params)
		case "upcoming":
			return application.GetUpcomingInvoice(ctx, cmd.SubscriptionID)
		}
	case "credit-notes":
		switch cmd.Operation {
		case "list":
			return application.ListCreditNotes(ctx, app.CatalogListParams{
				Limit:         cmd.Globals.Limit,
				Cursor:        cmd.Globals.Cursor,
				CreatedAfter:  cmd.CreatedAfter,
				CreatedBefore: cmd.CreatedBefore,
			})
		case "get":
			return application.GetCreditNote(ctx, cmd.ID)
		}
	}
	return output.Error(fmt.Errorf("unknown %s subcommand: %s", cmd.Resource, cmd.Operation), "usage_error", nil)
}

func executeEvents(ctx context.Context, application *app.App, cmd cli.Command) string {
	switch cmd.Operation {
	case "search":
		ids, err := eventIDs(cmd)
		if err != nil {
			return output.Error(err, "usage_error", map[string]interface{}{"resource": "events", "operation": "search"})
		}
		return application.SearchEvents(ctx, app.EventSearchParams{
			IDs:  ids,
			From: cmd.From,
			To:   cmd.To,
		})
	case "volume":
		return application.ListEventVolume(ctx, app.EventVolumeParams{
			Limit:  cmd.Globals.Limit,
			Cursor: cmd.Globals.Cursor,
			From:   cmd.From,
			To:     cmd.To,
		})
	case "backfills-list":
		return application.ListEventBackfills(ctx, app.PageParams{Limit: cmd.Globals.Limit, Cursor: cmd.Globals.Cursor})
	case "backfills-get":
		return application.GetEventBackfill(ctx, cmd.ID)
	default:
		return output.Error(fmt.Errorf("unknown events subcommand: %s", cmd.Operation), "usage_error", nil)
	}
}

func eventIDs(cmd cli.Command) ([]string, error) {
	ids := append([]string{}, cmd.IDs...)
	if cmd.IDsFile == "" {
		return ids, nil
	}

	content, err := os.ReadFile(cmd.IDsFile)
	if err != nil {
		return nil, fmt.Errorf("read --ids-file: %w", err)
	}
	for _, field := range strings.FieldsFunc(string(content), func(r rune) bool {
		return r == '\n' || r == '\r' || r == ','
	}) {
		id := strings.TrimSpace(field)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func executeAlerts(ctx context.Context, application *app.App, cmd cli.Command) string {
	switch cmd.Operation {
	case "list":
		return application.ListAlerts(ctx, app.AlertListParams{
			Limit:              cmd.Globals.Limit,
			Cursor:             cmd.Globals.Cursor,
			CreatedAfter:       cmd.CreatedAfter,
			CreatedBefore:      cmd.CreatedBefore,
			CustomerID:         cmd.CustomerID,
			ExternalCustomerID: cmd.ExternalCustomerID,
			SubscriptionID:     cmd.SubscriptionID,
		})
	case "get":
		return application.GetAlert(ctx, cmd.ID)
	default:
		return output.Error(fmt.Errorf("unknown alerts subcommand: %s", cmd.Operation), "usage_error", nil)
	}
}
