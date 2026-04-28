# orb Summary

`orb` is a Go CLI for Orb billing forensics.

Current implemented surface:

- `orb version`
- `orb ping`
- `orb customers list|get|costs|credits|credit-ledger`
- `orb subscriptions list|get|usage|costs|schedule`
- `orb plans list|get`
- `orb prices list|get`
- `orb metrics list|get`
- `orb invoices list|get|summary|upcoming`
- `orb credit-notes list|get`
- `orb events search|volume`
- `orb events backfills list|get`
- `orb alerts list|get`

Design goals:

- Read-heavy investigation workflows.
- JSON-first output for humans, scripts, Codex, and Claude.
- No real Orb API keys in tests.
- Release as `orb` through GitHub Releases and `vicentereig/homebrew-tap`.
