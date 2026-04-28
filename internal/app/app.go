package app

import (
	"context"
	"errors"
	"time"

	"github.com/vicentereig/orb/internal/output"
)

type OrbClient interface {
	Ping(ctx context.Context) (interface{}, error)
	ListCustomers(ctx context.Context, params CustomerListParams) (Page, error)
	GetCustomer(ctx context.Context, identity CustomerIdentity) (interface{}, error)
	ListCustomerCosts(ctx context.Context, identity CustomerIdentity, params TimeframeParams) (interface{}, error)
	ListCustomerCredits(ctx context.Context, identity CustomerIdentity, params CustomerCreditsParams) (Page, error)
	ListCustomerCreditLedger(ctx context.Context, identity CustomerIdentity, params CustomerCreditLedgerParams) (Page, error)
}

type Page struct {
	Data       interface{}
	Count      int
	NextCursor string
	HasMore    bool
}

type CustomerIdentity struct {
	ID         string
	ExternalID string
}

type CustomerListParams struct {
	Limit         int
	Cursor        string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
}

type TimeframeParams struct {
	From     *time.Time
	To       *time.Time
	Currency string
	ViewMode string
}

type CustomerCreditsParams struct {
	Limit            int
	Cursor           string
	Currency         string
	IncludeAllBlocks bool
}

type CustomerCreditLedgerParams struct {
	Limit       int
	Cursor      string
	Currency    string
	EntryType   string
	EntryStatus string
}

type App struct {
	client  OrbClient
	version string
}

func New(client OrbClient, version string) *App {
	return &App{client: client, version: version}
}

func (a *App) Version() string {
	return output.Success(
		map[string]string{"version": a.version},
		map[string]string{"resource": "system", "operation": "version"},
	)
}

func (a *App) Ping(ctx context.Context) string {
	res, err := a.client.Ping(ctx)
	if err != nil {
		return output.Error(err, "api_error", map[string]string{"resource": "system", "operation": "ping"})
	}
	return output.Success(res, map[string]string{"resource": "system", "operation": "ping"})
}

func (a *App) ListCustomers(ctx context.Context, params CustomerListParams) string {
	page, err := a.client.ListCustomers(ctx, params)
	if err != nil {
		return output.Error(err, "api_error", meta("customers", "list"))
	}
	return output.Success(page.Data, pageMeta("customers", "list", page))
}

func (a *App) GetCustomer(ctx context.Context, identity CustomerIdentity) string {
	if err := requireOneIdentity(identity); err != nil {
		return output.Error(err, "usage_error", meta("customers", "get"))
	}
	res, err := a.client.GetCustomer(ctx, identity)
	if err != nil {
		return output.Error(err, "api_error", meta("customers", "get"))
	}
	return output.Success(res, meta("customers", "get"))
}

func (a *App) ListCustomerCosts(ctx context.Context, identity CustomerIdentity, params TimeframeParams) string {
	if err := requireOneIdentity(identity); err != nil {
		return output.Error(err, "usage_error", meta("customers", "costs"))
	}
	res, err := a.client.ListCustomerCosts(ctx, identity, params)
	if err != nil {
		return output.Error(err, "api_error", meta("customers", "costs"))
	}
	return output.Success(res, meta("customers", "costs"))
}

func (a *App) ListCustomerCredits(ctx context.Context, identity CustomerIdentity, params CustomerCreditsParams) string {
	if err := requireOneIdentity(identity); err != nil {
		return output.Error(err, "usage_error", meta("customers", "credits"))
	}
	page, err := a.client.ListCustomerCredits(ctx, identity, params)
	if err != nil {
		return output.Error(err, "api_error", meta("customers", "credits"))
	}
	return output.Success(page.Data, pageMeta("customers", "credits", page))
}

func (a *App) ListCustomerCreditLedger(ctx context.Context, identity CustomerIdentity, params CustomerCreditLedgerParams) string {
	if err := requireOneIdentity(identity); err != nil {
		return output.Error(err, "usage_error", meta("customers", "credit-ledger"))
	}
	page, err := a.client.ListCustomerCreditLedger(ctx, identity, params)
	if err != nil {
		return output.Error(err, "api_error", meta("customers", "credit-ledger"))
	}
	return output.Success(page.Data, pageMeta("customers", "credit-ledger", page))
}

func requireOneIdentity(identity CustomerIdentity) error {
	if identity.ID == "" && identity.ExternalID == "" {
		return errors.New("--id or --external-id required")
	}
	if identity.ID != "" && identity.ExternalID != "" {
		return errors.New("--id and --external-id are mutually exclusive")
	}
	return nil
}

func meta(resource, operation string) map[string]interface{} {
	return map[string]interface{}{
		"resource":  resource,
		"operation": operation,
	}
}

func pageMeta(resource, operation string, page Page) map[string]interface{} {
	m := meta(resource, operation)
	m["count"] = page.Count
	m["has_more"] = page.HasMore
	if page.NextCursor != "" {
		m["next_cursor"] = page.NextCursor
	}
	return m
}
