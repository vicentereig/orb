package orbclient

import (
	"context"

	orbsdk "github.com/orbcorp/orb-go"
	"github.com/orbcorp/orb-go/option"
	"github.com/orbcorp/orb-go/packages/pagination"
	"github.com/vicentereig/orb/internal/app"
)

type Client struct {
	orb *orbsdk.Client
}

func New(apiKey, baseURL string) *Client {
	opts := []option.RequestOption{}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &Client{orb: orbsdk.NewClient(opts...)}
}

func (c *Client) Ping(ctx context.Context) (interface{}, error) {
	return c.orb.TopLevel.Ping(ctx)
}

func (c *Client) ListCustomers(ctx context.Context, params app.CustomerListParams) (app.Page, error) {
	page, err := c.orb.Customers.List(ctx, customerListParams(params))
	if err != nil {
		return app.Page{}, err
	}
	return pageFrom(page), nil
}

func (c *Client) GetCustomer(ctx context.Context, identity app.CustomerIdentity) (interface{}, error) {
	if identity.ExternalID != "" {
		return c.orb.Customers.FetchByExternalID(ctx, identity.ExternalID)
	}
	return c.orb.Customers.Fetch(ctx, identity.ID)
}

func (c *Client) ListCustomerCosts(ctx context.Context, identity app.CustomerIdentity, params app.TimeframeParams) (interface{}, error) {
	if identity.ExternalID != "" {
		return c.orb.Customers.Costs.ListByExternalID(ctx, identity.ExternalID, customerCostByExternalIDParams(params))
	}
	return c.orb.Customers.Costs.List(ctx, identity.ID, customerCostParams(params))
}

func (c *Client) ListCustomerCredits(ctx context.Context, identity app.CustomerIdentity, params app.CustomerCreditsParams) (app.Page, error) {
	if identity.ExternalID != "" {
		page, err := c.orb.Customers.Credits.ListByExternalID(ctx, identity.ExternalID, customerCreditByExternalIDParams(params))
		if err != nil {
			return app.Page{}, err
		}
		return pageFrom(page), nil
	}
	page, err := c.orb.Customers.Credits.List(ctx, identity.ID, customerCreditParams(params))
	if err != nil {
		return app.Page{}, err
	}
	return pageFrom(page), nil
}

func (c *Client) ListCustomerCreditLedger(ctx context.Context, identity app.CustomerIdentity, params app.CustomerCreditLedgerParams) (app.Page, error) {
	if identity.ExternalID != "" {
		page, err := c.orb.Customers.Credits.Ledger.ListByExternalID(ctx, identity.ExternalID, customerCreditLedgerByExternalIDParams(params))
		if err != nil {
			return app.Page{}, err
		}
		return pageFrom(page), nil
	}
	page, err := c.orb.Customers.Credits.Ledger.List(ctx, identity.ID, customerCreditLedgerParams(params))
	if err != nil {
		return app.Page{}, err
	}
	return pageFrom(page), nil
}

func customerListParams(params app.CustomerListParams) orbsdk.CustomerListParams {
	var q orbsdk.CustomerListParams
	if params.Limit > 0 {
		q.Limit = orbsdk.F(int64(params.Limit))
	}
	if params.Cursor != "" {
		q.Cursor = orbsdk.F(params.Cursor)
	}
	if params.CreatedAfter != nil {
		q.CreatedAtGte = orbsdk.F(*params.CreatedAfter)
	}
	if params.CreatedBefore != nil {
		q.CreatedAtLte = orbsdk.F(*params.CreatedBefore)
	}
	return q
}

func customerCostParams(params app.TimeframeParams) orbsdk.CustomerCostListParams {
	var q orbsdk.CustomerCostListParams
	if params.Currency != "" {
		q.Currency = orbsdk.F(params.Currency)
	}
	if params.From != nil {
		q.TimeframeStart = orbsdk.F(*params.From)
	}
	if params.To != nil {
		q.TimeframeEnd = orbsdk.F(*params.To)
	}
	if params.ViewMode != "" {
		q.ViewMode = orbsdk.F(orbsdk.CustomerCostListParamsViewMode(params.ViewMode))
	}
	return q
}

func customerCostByExternalIDParams(params app.TimeframeParams) orbsdk.CustomerCostListByExternalIDParams {
	var q orbsdk.CustomerCostListByExternalIDParams
	if params.Currency != "" {
		q.Currency = orbsdk.F(params.Currency)
	}
	if params.From != nil {
		q.TimeframeStart = orbsdk.F(*params.From)
	}
	if params.To != nil {
		q.TimeframeEnd = orbsdk.F(*params.To)
	}
	if params.ViewMode != "" {
		q.ViewMode = orbsdk.F(orbsdk.CustomerCostListByExternalIDParamsViewMode(params.ViewMode))
	}
	return q
}

func customerCreditParams(params app.CustomerCreditsParams) orbsdk.CustomerCreditListParams {
	var q orbsdk.CustomerCreditListParams
	if params.Limit > 0 {
		q.Limit = orbsdk.F(int64(params.Limit))
	}
	if params.Cursor != "" {
		q.Cursor = orbsdk.F(params.Cursor)
	}
	if params.Currency != "" {
		q.Currency = orbsdk.F(params.Currency)
	}
	if params.IncludeAllBlocks {
		q.IncludeAllBlocks = orbsdk.F(true)
	}
	return q
}

func customerCreditByExternalIDParams(params app.CustomerCreditsParams) orbsdk.CustomerCreditListByExternalIDParams {
	var q orbsdk.CustomerCreditListByExternalIDParams
	if params.Limit > 0 {
		q.Limit = orbsdk.F(int64(params.Limit))
	}
	if params.Cursor != "" {
		q.Cursor = orbsdk.F(params.Cursor)
	}
	if params.Currency != "" {
		q.Currency = orbsdk.F(params.Currency)
	}
	if params.IncludeAllBlocks {
		q.IncludeAllBlocks = orbsdk.F(true)
	}
	return q
}

func customerCreditLedgerParams(params app.CustomerCreditLedgerParams) orbsdk.CustomerCreditLedgerListParams {
	var q orbsdk.CustomerCreditLedgerListParams
	if params.Limit > 0 {
		q.Limit = orbsdk.F(int64(params.Limit))
	}
	if params.Cursor != "" {
		q.Cursor = orbsdk.F(params.Cursor)
	}
	if params.Currency != "" {
		q.Currency = orbsdk.F(params.Currency)
	}
	if params.EntryType != "" {
		q.EntryType = orbsdk.F(orbsdk.CustomerCreditLedgerListParamsEntryType(params.EntryType))
	}
	if params.EntryStatus != "" {
		q.EntryStatus = orbsdk.F(orbsdk.CustomerCreditLedgerListParamsEntryStatus(params.EntryStatus))
	}
	return q
}

func customerCreditLedgerByExternalIDParams(params app.CustomerCreditLedgerParams) orbsdk.CustomerCreditLedgerListByExternalIDParams {
	var q orbsdk.CustomerCreditLedgerListByExternalIDParams
	if params.Limit > 0 {
		q.Limit = orbsdk.F(int64(params.Limit))
	}
	if params.Cursor != "" {
		q.Cursor = orbsdk.F(params.Cursor)
	}
	if params.Currency != "" {
		q.Currency = orbsdk.F(params.Currency)
	}
	if params.EntryType != "" {
		q.EntryType = orbsdk.F(orbsdk.CustomerCreditLedgerListByExternalIDParamsEntryType(params.EntryType))
	}
	if params.EntryStatus != "" {
		q.EntryStatus = orbsdk.F(orbsdk.CustomerCreditLedgerListByExternalIDParamsEntryStatus(params.EntryStatus))
	}
	return q
}

func pageFrom[T any](page *pagination.Page[T]) app.Page {
	return app.Page{
		Data:       page.Data,
		Count:      len(page.Data),
		NextCursor: page.PaginationMetadata.NextCursor,
		HasMore:    page.PaginationMetadata.HasMore,
	}
}
