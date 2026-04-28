# orb CLI Plan

Date: 2026-04-28

## Goal

Build a Go CLI for read-heavy Orb investigations at work. The CLI should make it fast for a human or an agent such as Codex/Claude to inspect customers, subscriptions, plans, prices, usage, costs, invoices, credits, metrics, alerts, and known usage events.

The implementation should follow the practical shape of `whatsapp-cli`: single Go binary, stable JSON output, simple docs, GitHub Actions CI, tagged release workflow, and Homebrew distribution.

## Non-goals

- No real Orb API key in tests.
- No end-to-end tests against Orb.
- No local API key storage in the MVP.
- No mutating commands in the MVP, even when Orb exposes audit-safe mutations such as event amendment or deprecation.
- No local SQLite cache in the MVP. Add export/cache later only if repeated forensic workflows need it.

## Working Decisions

- Repository/module: `github.com/vicentereig/orb`.
- Binary and installed command name: `orb`.
- Go install path: `go install github.com/vicentereig/orb/cmd/orb@latest`, which installs a binary named `orb`.
- SDK: `github.com/orbcorp/orb-go`, pinned in `go.mod`.
- Go version: follow `whatsapp-cli` unless there is a reason not to. `orb-go` only requires Go 1.22+, so the implementation can stay compatible if desired.
- Auth: use `ORB_API_KEY` and optionally `ORB_BASE_URL`, matching `orb-go` defaults. Support explicit `--api-key` and `--base-url` only for one-off use, with redaction in all output.
- Output: JSON by default, no colors or progress output by default.
- Errors: non-zero exit code plus machine-readable JSON. Keep diagnostics out of stdout unless they are part of the JSON envelope.

## Project Shape

```text
withorb-cli/
  go.mod
  cmd/
    orb/
      main.go
  README.md
  QUICKSTART.md
  CHANGELOG.md
  SUMMARY.md
  docs/
    architecture.md
    RELEASE.md
    cli-api.md
    orb-api-map.md
  internal/
    app/
      app.go
      interfaces.go
      *_test.go
    cli/
      parser.go
      parser_test.go
    orbclient/
      client.go
      client_test.go
    output/
      output.go
      output_test.go
    params/
      time.go
      pagination.go
      *_test.go
  .github/
    workflows/
      ci.yml
      release.yml
```

This keeps the `whatsapp-cli` spirit but separates parsing from app logic early, because Orb has many similar resource commands and hand-parsing everything in `cmd/orb/main.go` will become brittle.

## Output Contract

All commands return the same envelope.

```json
{
  "success": true,
  "data": {},
  "meta": {
    "resource": "customers",
    "operation": "list",
    "count": 20,
    "next_cursor": "cursor-value"
  },
  "error": null
}
```

Failures use the same envelope.

```json
{
  "success": false,
  "data": null,
  "meta": {
    "resource": "customers",
    "operation": "get",
    "status_code": 404,
    "request_id": "req_..."
  },
  "error": {
    "message": "customer not found",
    "type": "api_error"
  }
}
```

For Codex/Claude friendliness:

- Stable snake_case field names.
- Compact JSON by default, `--pretty` for humans.
- Never print progress bars or incidental text unless `--verbose` is set.
- Redact `Authorization`, `ORB_API_KEY`, and custom secrets from errors and debug output.
- Include pagination metadata when the SDK response has it.

## Global Flags

```text
--api-key VALUE          Override ORB_API_KEY for this invocation
--base-url URL           Override ORB_BASE_URL for this invocation
--timeout DURATION       Request timeout, default 60s
--limit N                Page size for list commands, default 20
--cursor CURSOR          Pagination cursor
--all                    Auto-page through all results where safe
--pretty                 Pretty-print JSON
--raw                    Return SDK response shape with minimal normalization
--verbose                Include request diagnostics, with secrets redacted
```

`--all` should have a guardrail such as `--max-pages` before it ships, so an agent cannot accidentally pull an unbounded account export.

## MVP Commands

```text
orb version
orb ping

orb customers list [--created-after TIME] [--created-before TIME]
orb customers get --id CUSTOMER_ID
orb customers get --external-id EXTERNAL_CUSTOMER_ID
orb customers costs --id CUSTOMER_ID --from TIME --to TIME [--currency USD] [--view-mode cumulative|periodic]
orb customers costs --external-id EXTERNAL_CUSTOMER_ID --from TIME --to TIME
orb customers credits --id CUSTOMER_ID [--currency USD] [--include-all-blocks]
orb customers credits --external-id EXTERNAL_CUSTOMER_ID [--currency USD] [--include-all-blocks]
orb customers credit-ledger --id CUSTOMER_ID [--entry-type TYPE] [--entry-status STATUS]
orb customers credit-ledger --external-id EXTERNAL_CUSTOMER_ID [--entry-type TYPE] [--entry-status STATUS]

orb subscriptions list [--customer-id ID] [--external-customer-id ID] [--plan-id ID] [--external-plan-id ID] [--status active|ended|upcoming]
orb subscriptions get --id SUBSCRIPTION_ID
orb subscriptions usage --id SUBSCRIPTION_ID --from TIME --to TIME [--granularity day] [--group-by KEY] [--billable-metric-id ID]
orb subscriptions costs --id SUBSCRIPTION_ID --from TIME --to TIME [--currency USD] [--view-mode cumulative|periodic]
orb subscriptions schedule --id SUBSCRIPTION_ID

orb plans list [--status active|archived|draft]
orb plans get --id PLAN_ID
orb plans get --external-id EXTERNAL_PLAN_ID

orb prices list
orb prices get --id PRICE_ID
orb prices get --external-id EXTERNAL_PRICE_ID

orb metrics list
orb metrics get --id METRIC_ID

orb invoices list [--customer-id ID] [--external-customer-id ID] [--subscription-id ID] [--status STATUS] [--invoice-date-after TIME] [--invoice-date-before TIME]
orb invoices get --id INVOICE_ID
orb invoices summary [same filters as list]
orb invoices upcoming --subscription-id SUBSCRIPTION_ID

orb events search --id EVENT_ID [--id EVENT_ID ...] [--from TIME] [--to TIME]
orb events search --ids-file PATH [--from TIME] [--to TIME]
orb events volume --from TIME [--to TIME]
orb events backfills list
orb events backfills get --id BACKFILL_ID

orb alerts list [--customer-id ID] [--external-customer-id ID] [--subscription-id ID]
orb alerts get --id ALERT_ID

orb credit-notes list [--created-after TIME] [--created-before TIME]
orb credit-notes get --id CREDIT_NOTE_ID
```

Important event-search constraint: `orb-go` currently exposes event search around explicit event IDs, with optional timeframe bounds. It is not a general "find events by customer/event_name/property" API in the SDK surface I inspected. The CLI should document this clearly and avoid pretending it can enumerate arbitrary raw events unless a future Orb endpoint supports it.

## Phase 2 Forensic Commands

These are composite commands built on top of the MVP primitives.

```text
orb customers dossier --id CUSTOMER_ID --from TIME --to TIME
orb customers dossier --external-id EXTERNAL_CUSTOMER_ID --from TIME --to TIME
orb subscriptions dossier --id SUBSCRIPTION_ID --from TIME --to TIME
orb prices explain --id PRICE_ID
orb plans explain --id PLAN_ID
orb raw get /path [--query key=value ...]
```

`customers dossier` should gather customer details, subscriptions, usage/costs by subscription, customer-level costs, credits, credit ledger entries, invoices, alerts, and optionally known events via `--event-id` or `--ids-file`.

`raw get` is a read-only escape hatch using `client.Get`/`client.Execute` for SDK gaps. It should be excluded from the MVP until the normal commands are solid.

Credit note list filtering is limited in the current SDK to creation time, cursor, and limit. Customer-specific credit note investigation should go through invoice/customer context until Orb exposes direct customer filters for credit notes.

## TDD Strategy

Tests should be written before implementation for each command family.

- Parser tests: argv to command structs, including subcommand flag positioning and validation errors.
- App tests: fake Orb client interface, asserting output behavior instead of SDK internals.
- Output tests: success/error envelope, redaction, pretty output, API error conversion.
- Param tests: time parsing, `--from`/`--to`, cursor/limit, repeated flags, comma-separated lists.
- Adapter tests: use `httptest.Server` and `orb-go` configured with fake `ORB_API_KEY` and fake `ORB_BASE_URL` only where we need confidence that SDK params serialize as expected. These are local contract tests, not Orb e2e tests.
- No tests should require network access, real credentials, or recorded production payloads.

Recommended red-green order:

1. `version`, `ping`, output envelope.
2. Global parser and auth option construction.
3. Customers list/get/external-get.
4. Pagination helper and `--all` guardrails.
5. Subscriptions list/get/usage/costs/schedule.
6. Plans, prices, metrics.
7. Invoices and credit notes.
8. Events search/volume/backfills.
9. Alerts.
10. Composite dossier commands.

## Implementation Phases

1. Scaffold the Go module and command entrypoint.
   - Create `cmd/orb/main.go`, `internal/output`, `internal/cli`, and the initial `version` command.
   - Add tests first for output envelopes, parser behavior, and version output.

2. Build the read-only Orb API primitives.
   - Implement customers, subscriptions, catalog, billing, events, and alerts in the red-green order above.
   - Keep all API access behind interfaces so tests use fakes or local `httptest` servers, never a real `ORB_API_KEY`.

3. Put together GitHub Actions and release automation from `whatsapp-cli`.
   - Add `.github/workflows/ci.yml` with checkout, `actions/setup-go` using `go.mod`, and `go test ./...`.
   - Add `.github/workflows/release.yml` triggered by `v*` tags.
   - Run tests before release builds.
   - Build `orb` artifacts for Linux `amd64`/`arm64`, macOS `amd64`/`arm64`, and Windows `amd64`.
   - Package artifacts as `orb-<os>-<arch>.tar.gz` and `orb-windows-amd64.zip`.
   - Generate `checksums.txt`.
   - Sign artifacts and checksums with Sigstore cosign using GitHub OIDC.
   - Publish artifacts to the GitHub Release for the tag.
   - Update `vicentereig/homebrew-tap` after publish, using `HOMEBREW_TAP_TOKEN`, with a formula named `orb.rb` unless the tap already needs a different convention.

4. Write release and installation docs.
   - Add `docs/RELEASE.md` modeled on `whatsapp-cli`, but simplified because `orb` should be CGO-free.
  - Document `go install github.com/vicentereig/orb/cmd/orb@latest`.
   - Document Homebrew install as `brew install vicentereig/tap/orb`.
   - Include checksum and cosign verification examples.

5. Add composite forensic commands.
   - Implement `customers dossier`, `subscriptions dossier`, and explain commands once the primitive commands and release pipeline are stable.

## CI And Release

Follow `whatsapp-cli`:

- `.github/workflows/ci.yml`: checkout, setup Go from `go.mod`, `go test ./...`.
- Add `go vet ./...` once the initial surface settles.
- `.github/workflows/release.yml`: tagged `v*` releases, run tests, build Linux/macOS/Windows `orb` artifacts, checksums, Sigstore cosign signatures, GitHub Release upload, and `vicentereig/homebrew-tap` update.
- `docs/RELEASE.md`: semver tagging, checksums, cosign verification, Homebrew update behavior.

Because this project does not use CGO, release builds should be simpler than `whatsapp-cli`; no SQLite or cross-compiler toolchain should be needed.

## Documentation

Docs should be written for humans and agents:

- `README.md`: complete reference with examples, JSON schema, installation, auth, troubleshooting.
- `QUICKSTART.md`: first useful commands with fake-safe examples.
- `docs/architecture.md`: command flow, test strategy, SDK wrapper boundary.
- `docs/cli-api.md`: command reference and examples.
- `docs/orb-api-map.md`: SDK methods backing each command.
- `SUMMARY.md`: compact overview for LLM context.

## Open Questions

- Should the first release be strictly read-only, or include event `deprecate`/`amend` behind a loud `--confirm` gate?
- Do you want profiles for multiple Orb accounts, or is environment-variable auth enough?
- Should `dossier` output include only raw fetched objects, or also computed summaries such as totals by invoice status and metric?
