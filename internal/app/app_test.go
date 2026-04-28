package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type fakeClient struct {
	pingResponse   interface{}
	pingErr        error
	pageResponse   Page
	objectResponse interface{}
	err            error
	lastCustomer   CustomerIdentity
}

func (f fakeClient) Ping(ctx context.Context) (interface{}, error) {
	if f.pingErr != nil {
		return nil, f.pingErr
	}
	return f.pingResponse, nil
}

func (f fakeClient) ListCustomers(ctx context.Context, params CustomerListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetCustomer(ctx context.Context, identity CustomerIdentity) (interface{}, error) {
	f.lastCustomer = identity
	return f.objectResponse, f.err
}

func (f fakeClient) ListCustomerCosts(ctx context.Context, identity CustomerIdentity, params TimeframeParams) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListCustomerCredits(ctx context.Context, identity CustomerIdentity, params CustomerCreditsParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) ListCustomerCreditLedger(ctx context.Context, identity CustomerIdentity, params CustomerCreditLedgerParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) ListSubscriptions(ctx context.Context, params SubscriptionListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetSubscription(ctx context.Context, subscriptionID string) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListSubscriptionUsage(ctx context.Context, subscriptionID string, params SubscriptionUsageParams) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListSubscriptionCosts(ctx context.Context, subscriptionID string, params TimeframeParams) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListSubscriptionSchedule(ctx context.Context, subscriptionID string, params SubscriptionScheduleParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) ListPlans(ctx context.Context, params CatalogListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetPlan(ctx context.Context, identity ResourceIdentity) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListPrices(ctx context.Context, params PageParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetPrice(ctx context.Context, identity ResourceIdentity) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListMetrics(ctx context.Context, params CatalogListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetMetric(ctx context.Context, metricID string) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListInvoices(ctx context.Context, params InvoiceListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetInvoice(ctx context.Context, invoiceID string) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListInvoiceSummary(ctx context.Context, params InvoiceListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetUpcomingInvoice(ctx context.Context, subscriptionID string) (interface{}, error) {
	return f.objectResponse, f.err
}

func (f fakeClient) ListCreditNotes(ctx context.Context, params CatalogListParams) (Page, error) {
	return f.pageResponse, f.err
}

func (f fakeClient) GetCreditNote(ctx context.Context, creditNoteID string) (interface{}, error) {
	return f.objectResponse, f.err
}

func TestVersion(t *testing.T) {
	result := New(fakeClient{}, "v1.2.3").Version()

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data.Version != "v1.2.3" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestPingSuccess(t *testing.T) {
	result := New(fakeClient{pingResponse: map[string]string{"response": "pong"}}, "dev").Ping(context.Background())

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			Response string `json:"response"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data.Response != "pong" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestPingError(t *testing.T) {
	result := New(fakeClient{pingErr: errors.New("unauthorized")}, "dev").Ping(context.Background())

	var payload struct {
		Success bool `json:"success"`
		Error   struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Success || payload.Error.Message != "unauthorized" || payload.Error.Type != "api_error" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListCustomersSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{
		Data:       []map[string]string{{"id": "cus_123"}},
		Count:      1,
		NextCursor: "next",
		HasMore:    true,
	}}, "dev")

	result := app.ListCustomers(context.Background(), CustomerListParams{Limit: 20, Cursor: "start"})

	var payload struct {
		Success bool                   `json:"success"`
		Data    []map[string]string    `json:"data"`
		Meta    map[string]interface{} `json:"meta"`
		Error   map[string]interface{} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 || payload.Meta["next_cursor"] != "next" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetCustomerRequiresOneIdentity(t *testing.T) {
	app := New(fakeClient{}, "dev")

	result := app.GetCustomer(context.Background(), CustomerIdentity{})

	var payload struct {
		Success bool `json:"success"`
		Error   struct {
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Success || payload.Error.Type != "usage_error" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetCustomerSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"id": "cus_123"}}, "dev")

	result := app.GetCustomer(context.Background(), CustomerIdentity{ID: "cus_123"})

	var payload struct {
		Success bool              `json:"success"`
		Data    map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data["id"] != "cus_123" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestCustomerCostsSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"total": "42"}}, "dev")

	result := app.ListCustomerCosts(context.Background(), CustomerIdentity{ExternalID: "workspace_123"}, TimeframeParams{})

	var payload struct {
		Success bool              `json:"success"`
		Data    map[string]string `json:"data"`
		Meta    map[string]string `json:"meta"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Meta["operation"] != "costs" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestCustomerCreditsSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"block": "credits"}}, Count: 1}}, "dev")

	result := app.ListCustomerCredits(context.Background(), CustomerIdentity{ID: "cus_123"}, CustomerCreditsParams{Limit: 10})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestCustomerCreditLedgerSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"entry": "increment"}}, Count: 1}}, "dev")

	result := app.ListCustomerCreditLedger(context.Background(), CustomerIdentity{ID: "cus_123"}, CustomerCreditLedgerParams{EntryType: "increment"})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListSubscriptionsSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "sub_123"}}, Count: 1}}, "dev")

	result := app.ListSubscriptions(context.Background(), SubscriptionListParams{CustomerID: "cus_123", Status: "active"})

	var payload struct {
		Success bool                   `json:"success"`
		Data    []map[string]string    `json:"data"`
		Meta    map[string]interface{} `json:"meta"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 || payload.Meta["operation"] != "list" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetSubscriptionRequiresID(t *testing.T) {
	app := New(fakeClient{}, "dev")

	result := app.GetSubscription(context.Background(), "")

	var payload struct {
		Success bool `json:"success"`
		Error   struct {
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Success || payload.Error.Type != "usage_error" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetSubscriptionSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"id": "sub_123"}}, "dev")

	result := app.GetSubscription(context.Background(), "sub_123")

	var payload struct {
		Success bool              `json:"success"`
		Data    map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data["id"] != "sub_123" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestSubscriptionUsageSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"usage": "42"}}, "dev")

	result := app.ListSubscriptionUsage(context.Background(), "sub_123", SubscriptionUsageParams{GroupBy: "region"})

	var payload struct {
		Success bool              `json:"success"`
		Data    map[string]string `json:"data"`
		Meta    map[string]string `json:"meta"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Meta["operation"] != "usage" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestSubscriptionCostsSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"cost": "42"}}, "dev")

	result := app.ListSubscriptionCosts(context.Background(), "sub_123", TimeframeParams{Currency: "USD"})

	var payload struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestSubscriptionScheduleSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"phase": "next"}}, Count: 1}}, "dev")

	result := app.ListSubscriptionSchedule(context.Background(), "sub_123", SubscriptionScheduleParams{Limit: 10})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListPlansSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "plan_123"}}, Count: 1}}, "dev")

	result := app.ListPlans(context.Background(), CatalogListParams{Status: "active"})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetPlanSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"id": "plan_123"}}, "dev")

	result := app.GetPlan(context.Background(), ResourceIdentity{ID: "plan_123"})

	var payload struct {
		Success bool              `json:"success"`
		Data    map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data["id"] != "plan_123" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListPricesSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "price_123"}}, Count: 1}}, "dev")

	result := app.ListPrices(context.Background(), PageParams{Limit: 50})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetPriceSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"id": "price_123"}}, "dev")

	result := app.GetPrice(context.Background(), ResourceIdentity{ExternalID: "api_calls"})

	var payload struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListMetricsSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "metric_123"}}, Count: 1}}, "dev")

	result := app.ListMetrics(context.Background(), CatalogListParams{Limit: 20})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetMetricRequiresID(t *testing.T) {
	app := New(fakeClient{}, "dev")

	result := app.GetMetric(context.Background(), "")

	var payload struct {
		Success bool `json:"success"`
		Error   struct {
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Success || payload.Error.Type != "usage_error" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListInvoicesSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "inv_123"}}, Count: 1}}, "dev")

	result := app.ListInvoices(context.Background(), InvoiceListParams{CustomerID: "cus_123", Status: "issued"})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestGetInvoiceSuccess(t *testing.T) {
	app := New(fakeClient{objectResponse: map[string]string{"id": "inv_123"}}, "dev")

	result := app.GetInvoice(context.Background(), "inv_123")

	var payload struct {
		Success bool              `json:"success"`
		Data    map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data["id"] != "inv_123" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestInvoiceSummarySuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "summary_123"}}, Count: 1}}, "dev")

	result := app.ListInvoiceSummary(context.Background(), InvoiceListParams{ExternalCustomerID: "workspace_123"})

	var payload struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestUpcomingInvoiceRequiresSubscriptionID(t *testing.T) {
	app := New(fakeClient{}, "dev")

	result := app.GetUpcomingInvoice(context.Background(), "")

	var payload struct {
		Success bool `json:"success"`
		Error   struct {
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Success || payload.Error.Type != "usage_error" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestListCreditNotesSuccess(t *testing.T) {
	app := New(fakeClient{pageResponse: Page{Data: []map[string]string{{"id": "cn_123"}}, Count: 1}}, "dev")

	result := app.ListCreditNotes(context.Background(), CatalogListParams{Limit: 20})

	var payload struct {
		Success bool                `json:"success"`
		Data    []map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || len(payload.Data) != 1 {
		t.Fatalf("unexpected result: %s", result)
	}
}
