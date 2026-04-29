# orb

Chat with your revenue in Orb.

`orb` is a read-first CLI for investigating Orb customers, subscriptions, plans, prices, usage, invoices, credits, events, and related billing state. It is designed for humans and coding agents: stable JSON output, predictable flags, and no incidental terminal noise.

## Install

```sh
go install github.com/vicentereig/orb/cmd/orb@latest
```

After the first release:

```sh
brew install vicentereig/tap/orb
```

## Authentication

The CLI uses the same environment variables as `github.com/orbcorp/orb-go`.

```sh
export ORB_API_KEY="..."
orb ping
```

Optional:

```sh
export ORB_BASE_URL="https://api.withorb.com"
```

## Current Commands

```sh
orb help
orb help events
orb help examples
orb version
orb ping

orb customers list --limit 50
orb customers get --id cus_123
orb customers get --external-id workspace_123
orb customers costs --id cus_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z
orb customers credits --id cus_123 --include-all-blocks
orb customers credit-ledger --id cus_123 --entry-status committed

orb subscriptions list --customer-id cus_123 --status active
orb subscriptions get --id sub_123
orb subscriptions usage --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z
orb subscriptions costs --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z
orb subscriptions schedule --id sub_123

orb plans list
orb plans get --id plan_123
orb prices list
orb prices get --id price_123
orb metrics list
orb metrics get --id metric_123

orb invoices list --customer-id cus_123
orb invoices get --id inv_123
orb invoices summary --customer-id cus_123
orb invoices upcoming --subscription-id sub_123
orb credit-notes list
orb credit-notes get --id cn_123

orb events search --id event_id_1 --id event_id_2
orb events search --ids-file ./event_ids.txt
orb events volume --from 2026-04-01T00:00:00Z
orb events backfills list
orb events backfills get --id backfill_123

orb alerts list --customer-id cus_123
orb alerts get --id alert_123
```

Event search maps to Orb's explicit event ID search API. It is not a general event query by customer, name, or property.

## Agent Help

`orb help` returns structured JSON with topics, global flags, commands, examples, and notes. Use topic help when an agent needs command-specific required flags.

```sh
orb help --pretty
orb help customers --pretty
orb help events --pretty
```

Every command returns JSON:

```json
{
  "success": true,
  "data": {
    "version": "dev"
  },
  "meta": {
    "operation": "version",
    "resource": "system"
  },
  "error": null
}
```

## Development

```sh
go test ./...
go run ./cmd/orb version
```

The implementation is TDD-first and does not use a real Orb API key in tests. API-facing tests should use fakes or local `httptest` servers.
