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
	ListSubscriptions(ctx context.Context, params SubscriptionListParams) (Page, error)
	GetSubscription(ctx context.Context, subscriptionID string) (interface{}, error)
	ListSubscriptionUsage(ctx context.Context, subscriptionID string, params SubscriptionUsageParams) (interface{}, error)
	ListSubscriptionCosts(ctx context.Context, subscriptionID string, params TimeframeParams) (interface{}, error)
	ListSubscriptionSchedule(ctx context.Context, subscriptionID string, params SubscriptionScheduleParams) (Page, error)
	ListPlans(ctx context.Context, params CatalogListParams) (Page, error)
	GetPlan(ctx context.Context, identity ResourceIdentity) (interface{}, error)
	ListPrices(ctx context.Context, params PageParams) (Page, error)
	GetPrice(ctx context.Context, identity ResourceIdentity) (interface{}, error)
	ListMetrics(ctx context.Context, params CatalogListParams) (Page, error)
	GetMetric(ctx context.Context, metricID string) (interface{}, error)
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

type ResourceIdentity struct {
	ID         string
	ExternalID string
}

type PageParams struct {
	Limit  int
	Cursor string
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

type SubscriptionListParams struct {
	Limit              int
	Cursor             string
	CreatedAfter       *time.Time
	CreatedBefore      *time.Time
	CustomerID         string
	ExternalCustomerID string
	PlanID             string
	ExternalPlanID     string
	Status             string
}

type SubscriptionUsageParams struct {
	From                 *time.Time
	To                   *time.Time
	BillableMetricID     string
	FirstDimensionKey    string
	FirstDimensionValue  string
	Granularity          string
	GroupBy              string
	SecondDimensionKey   string
	SecondDimensionValue string
	ViewMode             string
}

type SubscriptionScheduleParams struct {
	Limit       int
	Cursor      string
	StartAfter  *time.Time
	StartBefore *time.Time
}

type CatalogListParams struct {
	Limit         int
	Cursor        string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	Status        string
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

func (a *App) ListSubscriptions(ctx context.Context, params SubscriptionListParams) string {
	page, err := a.client.ListSubscriptions(ctx, params)
	if err != nil {
		return output.Error(err, "api_error", meta("subscriptions", "list"))
	}
	return output.Success(page.Data, pageMeta("subscriptions", "list", page))
}

func (a *App) GetSubscription(ctx context.Context, subscriptionID string) string {
	if subscriptionID == "" {
		return output.Error(errors.New("--id required"), "usage_error", meta("subscriptions", "get"))
	}
	res, err := a.client.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return output.Error(err, "api_error", meta("subscriptions", "get"))
	}
	return output.Success(res, meta("subscriptions", "get"))
}

func (a *App) ListSubscriptionUsage(ctx context.Context, subscriptionID string, params SubscriptionUsageParams) string {
	if subscriptionID == "" {
		return output.Error(errors.New("--id required"), "usage_error", meta("subscriptions", "usage"))
	}
	res, err := a.client.ListSubscriptionUsage(ctx, subscriptionID, params)
	if err != nil {
		return output.Error(err, "api_error", meta("subscriptions", "usage"))
	}
	return output.Success(res, meta("subscriptions", "usage"))
}

func (a *App) ListSubscriptionCosts(ctx context.Context, subscriptionID string, params TimeframeParams) string {
	if subscriptionID == "" {
		return output.Error(errors.New("--id required"), "usage_error", meta("subscriptions", "costs"))
	}
	res, err := a.client.ListSubscriptionCosts(ctx, subscriptionID, params)
	if err != nil {
		return output.Error(err, "api_error", meta("subscriptions", "costs"))
	}
	return output.Success(res, meta("subscriptions", "costs"))
}

func (a *App) ListSubscriptionSchedule(ctx context.Context, subscriptionID string, params SubscriptionScheduleParams) string {
	if subscriptionID == "" {
		return output.Error(errors.New("--id required"), "usage_error", meta("subscriptions", "schedule"))
	}
	page, err := a.client.ListSubscriptionSchedule(ctx, subscriptionID, params)
	if err != nil {
		return output.Error(err, "api_error", meta("subscriptions", "schedule"))
	}
	return output.Success(page.Data, pageMeta("subscriptions", "schedule", page))
}

func (a *App) ListPlans(ctx context.Context, params CatalogListParams) string {
	page, err := a.client.ListPlans(ctx, params)
	if err != nil {
		return output.Error(err, "api_error", meta("plans", "list"))
	}
	return output.Success(page.Data, pageMeta("plans", "list", page))
}

func (a *App) GetPlan(ctx context.Context, identity ResourceIdentity) string {
	if err := requireOneResourceIdentity(identity); err != nil {
		return output.Error(err, "usage_error", meta("plans", "get"))
	}
	res, err := a.client.GetPlan(ctx, identity)
	if err != nil {
		return output.Error(err, "api_error", meta("plans", "get"))
	}
	return output.Success(res, meta("plans", "get"))
}

func (a *App) ListPrices(ctx context.Context, params PageParams) string {
	page, err := a.client.ListPrices(ctx, params)
	if err != nil {
		return output.Error(err, "api_error", meta("prices", "list"))
	}
	return output.Success(page.Data, pageMeta("prices", "list", page))
}

func (a *App) GetPrice(ctx context.Context, identity ResourceIdentity) string {
	if err := requireOneResourceIdentity(identity); err != nil {
		return output.Error(err, "usage_error", meta("prices", "get"))
	}
	res, err := a.client.GetPrice(ctx, identity)
	if err != nil {
		return output.Error(err, "api_error", meta("prices", "get"))
	}
	return output.Success(res, meta("prices", "get"))
}

func (a *App) ListMetrics(ctx context.Context, params CatalogListParams) string {
	page, err := a.client.ListMetrics(ctx, params)
	if err != nil {
		return output.Error(err, "api_error", meta("metrics", "list"))
	}
	return output.Success(page.Data, pageMeta("metrics", "list", page))
}

func (a *App) GetMetric(ctx context.Context, metricID string) string {
	if metricID == "" {
		return output.Error(errors.New("--id required"), "usage_error", meta("metrics", "get"))
	}
	res, err := a.client.GetMetric(ctx, metricID)
	if err != nil {
		return output.Error(err, "api_error", meta("metrics", "get"))
	}
	return output.Success(res, meta("metrics", "get"))
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

func requireOneResourceIdentity(identity ResourceIdentity) error {
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
