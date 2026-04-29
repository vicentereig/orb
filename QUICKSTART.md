# Quick Start

## Build

```sh
go build -o orb ./cmd/orb
```

## Check Version

```sh
./orb version
```

## Check Orb Credentials

```sh
export ORB_API_KEY="..."
./orb ping
```

## Pretty JSON

```sh
./orb --pretty version
```

## Learn Commands

```sh
./orb help --pretty
./orb help events --pretty
./orb help examples --pretty
```

Help output is JSON by default so humans, scripts, Codex, and Claude can inspect the same command metadata.

## Read Billing State

```sh
./orb customers get --id cus_123
./orb subscriptions list --customer-id cus_123
./orb invoices list --customer-id cus_123
./orb alerts list --customer-id cus_123
```

## Investigate Usage

```sh
./orb subscriptions usage --id sub_123 --from 2026-04-01T00:00:00Z --to 2026-05-01T00:00:00Z
./orb events volume --from 2026-04-01T00:00:00Z --to 2026-04-02T00:00:00Z
./orb events search --id event_id_1 --id event_id_2
```

`events search` requires explicit event IDs. For larger forensic lookups, put one event ID per line in a file and use:

```sh
./orb events search --ids-file ./event_ids.txt
```

## Tests

```sh
go test ./...
```
