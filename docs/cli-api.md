# CLI API Design

Date: 2026-04-28

## Principles

- Read-first. Forensic workflows should not mutate production billing state by accident.
- JSON-first. Every command returns one machine-readable envelope.
- Predictable flags. Resource identifiers use `--id`, external identifiers use `--external-id`, time ranges use `--from` and `--to`.
- Pagination-aware. List commands expose `--limit`, `--cursor`, and later `--all --max-pages`.
- SDK-shaped, not SDK-leaky. Commands should map cleanly to Orb concepts, while hiding generated Go naming where it does not help the operator.

## Envelope

```json
{
  "success": true,
  "data": {},
  "meta": {
    "resource": "subscriptions",
    "operation": "usage",
    "next_cursor": null
  },
  "error": null
}
```

`data` should usually be the Orb SDK response or the SDK page `data` array. `meta` carries pagination, request, and CLI context.

## Core Commands

### Health

```sh
orb version
orb ping
```

### Customers

```sh
orb customers list --limit 50
orb customers list --created-after 2026-01-01T00:00:00Z
orb customers get --id cus_123
orb customers get --external-id workspace_123
orb customers costs --id cus_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z --view-mode cumulative
orb customers credits --id cus_123 --include-all-blocks
orb customers credits --external-id workspace_123 --include-all-blocks
orb customers credit-ledger --id cus_123 --entry-status committed --entry-type increment
orb customers credit-ledger --external-id workspace_123 --entry-status committed
```

### Subscriptions

```sh
orb subscriptions list --customer-id cus_123 --status active
orb subscriptions list --external-customer-id workspace_123
orb subscriptions get --id sub_123
orb subscriptions usage --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z --granularity day
orb subscriptions usage --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z --group-by region --billable-metric-id bm_123
orb subscriptions costs --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z --currency USD
orb subscriptions schedule --id sub_123 --limit 100
```

### Catalog

```sh
orb plans list --status active
orb plans get --id plan_123
orb plans get --external-id pro_plan
orb prices list --limit 100
orb prices get --id price_123
orb prices get --external-id api_calls_price
orb metrics list
orb metrics get --id metric_123
```

### Billing

```sh
orb invoices list --customer-id cus_123 --status issued --invoice-date-after 2026-04-01T00:00:00Z
orb invoices list --external-customer-id workspace_123 --subscription-id sub_123
orb invoices get --id inv_123
orb invoices summary --customer-id cus_123
orb invoices upcoming --subscription-id sub_123
orb credit-notes list --created-after 2026-04-01T00:00:00Z
orb credit-notes get --id cn_123
```

### Events

```sh
orb events search --id event_id_1 --id event_id_2 --from 2026-04-01T00:00:00Z --to 2026-04-02T00:00:00Z
orb events search --ids-file ./event_ids.txt
orb events volume --from 2026-04-01T00:00:00Z --to 2026-04-02T00:00:00Z
orb events backfills list
orb events backfills get --id backfill_123
```

Current SDK limitation: event search requires explicit event IDs. Do not design flags such as `--customer-id`, `--event-name`, or `--property` for event search until Orb exposes that search shape.

### Alerts

```sh
orb alerts list
orb alerts list --customer-id cus_123
orb alerts get --id alert_123
```

## Composite Commands

Composite commands are phase 2, after the primitives are tested.

```sh
orb customers dossier --external-id workspace_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z
orb subscriptions dossier --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z
orb plans explain --id plan_123
orb prices explain --id price_123
```

`dossier` output should be a structured object with sections:

```json
{
  "customer": {},
  "subscriptions": [],
  "usage": [],
  "costs": [],
  "invoices": [],
  "credits": [],
  "credit_ledger": [],
  "alerts": [],
  "events": []
}
```

If event IDs are not provided, `events` should be omitted or include only event volume data, not an empty raw-event claim.

## Validation Rules

- `get` commands require exactly one identity flag.
- External ID flags must not be mixed with Orb ID flags on the same command.
- Time ranges require `--from` before `--to` when both are present.
- `--all` requires a maximum page limit before release.
- Mutating SDK methods are not wired in the MVP.
