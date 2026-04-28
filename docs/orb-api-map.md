# Orb SDK API Map

Date: 2026-04-28

This map is based on `github.com/orbcorp/orb-go` `main` as inspected on 2026-04-28.

## SDK Notes

- `orb.NewClient()` reads `ORB_API_KEY`, `ORB_WEBHOOK_SECRET`, and `ORB_BASE_URL` from the environment by default.
- `option.WithAPIKey` and `option.WithBaseURL` can override env values.
- Request params use `orb.F(...)` field wrappers, so CLI param builders should live in one small package and be well tested.
- Paginated list endpoints return `pagination.Page[T]` and support cursor/limit style pagination.
- The SDK has retry behavior for connection errors, 408, 409, 429, and 5xx by default. The CLI should still apply a context timeout.

## Read Command Mapping

| CLI area | SDK methods |
| --- | --- |
| `ping` | `client.TopLevel.Ping` |
| `customers list` | `client.Customers.List` |
| `customers get --id` | `client.Customers.Fetch` |
| `customers get --external-id` | `client.Customers.FetchByExternalID` |
| `customers costs` | `client.Customers.Costs.List`, `client.Customers.Costs.ListByExternalID` |
| `customers credits` | `client.Customers.Credits.List`, `client.Customers.Credits.ListByExternalID` |
| `customers credit-ledger` | `client.Customers.Credits.Ledger.List`, `client.Customers.Credits.Ledger.ListByExternalID` |
| `subscriptions list` | `client.Subscriptions.List` |
| `subscriptions get` | `client.Subscriptions.Fetch` |
| `subscriptions usage` | `client.Subscriptions.FetchUsage` |
| `subscriptions costs` | `client.Subscriptions.FetchCosts` |
| `subscriptions schedule` | `client.Subscriptions.FetchSchedule` |
| `plans list` | `client.Plans.List` |
| `plans get --id` | `client.Plans.Fetch` |
| `plans get --external-id` | `client.Plans.ExternalPlanID.Fetch` |
| `prices list` | `client.Prices.List` |
| `prices get --id` | `client.Prices.Fetch` |
| `prices get --external-id` | `client.Prices.ExternalPriceID.Fetch` |
| `metrics list` | `client.Metrics.List` |
| `metrics get` | `client.Metrics.Fetch` |
| `invoices list` | `client.Invoices.List` |
| `invoices get` | `client.Invoices.Fetch` |
| `invoices summary` | `client.Invoices.ListSummary` |
| `invoices upcoming` | `client.Invoices.FetchUpcoming` |
| `credit-notes list` | `client.CreditNotes.List` |
| `credit-notes get` | `client.CreditNotes.Fetch` |
| `events search` | `client.Events.Search` |
| `events volume` | `client.Events.Volume.List` |
| `events backfills list` | `client.Events.Backfills.List` |
| `events backfills get` | `client.Events.Backfills.Fetch` |
| `alerts list` | `client.Alerts.List` |
| `alerts get` | `client.Alerts.Get` |

## Useful Params To Support First

| Resource | Parameters |
| --- | --- |
| Customers list | `created_at[gt/gte/lt/lte]`, `cursor`, `limit` |
| Customer costs | `timeframe_start`, `timeframe_end`, `currency`, `view_mode` |
| Customer credits | `currency`, `cursor`, `effective_date[...]`, `include_all_blocks`, `limit` |
| Customer credit ledger | `created_at[...]`, `currency`, `cursor`, `entry_status`, `entry_type`, `limit`, `minimum_amount` |
| Subscriptions list | `created_at[...]`, `cursor`, `customer_id[]`, `external_customer_id[]`, `external_plan_id`, `limit`, `plan_id`, `status` |
| Subscription usage | `billable_metric_id`, dimensions, `granularity`, `group_by`, `timeframe_start`, `timeframe_end`, `view_mode` |
| Subscription costs | `currency`, `timeframe_start`, `timeframe_end`, `view_mode` |
| Subscription schedule | `cursor`, `limit`, `start_date[...]` |
| Plans list | `created_at[...]`, `cursor`, `limit`, `status` |
| Prices list | `cursor`, `limit` |
| Invoices list | `amount`, `amount[gt/lt]`, `cursor`, `customer_id`, `date_type`, due-date filters, `external_customer_id`, invoice-date filters, `is_recurring`, `limit`, `status[]`, `subscription_id` |
| Event search | `event_ids[]`, `timeframe_start`, `timeframe_end` |
| Event volume | required `timeframe_start`, optional `timeframe_end`, `cursor`, `limit` |
| Alerts list | `created_at[...]`, `cursor`, `customer_id`, `external_customer_id`, `limit`, `subscription_id` |
| Credit notes list | `created_at[...]`, `cursor`, `limit` |

## Mutating Methods Intentionally Excluded From MVP

- Customer create/update/delete/sync payment methods.
- Subscription create/update/cancel/schedule mutations.
- Invoice create/update/issue/pay/void/mark paid/delete line item.
- Event ingest/amend/deprecate.
- Event backfill create/close/revert.
- Alert create/update/enable/disable.
- Credit ledger entry/top-up mutations.

These can be added later behind explicit command names and confirmation gates if needed.
