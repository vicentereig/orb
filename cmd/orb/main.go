package main

import (
	"context"
	"fmt"
	"os"

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
  version   Print CLI version information
  ping      Check Orb API connectivity and credentials

Global options:
  --api-key VALUE     Override ORB_API_KEY for this invocation
  --base-url URL      Override ORB_BASE_URL for this invocation
  --timeout DURATION  Request timeout (default: 60s)
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
