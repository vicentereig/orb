# orb

Orb billing forensics from the terminal.

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
orb version
orb ping
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
