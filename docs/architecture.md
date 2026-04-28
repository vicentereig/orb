# Architecture

`orb` is a single-binary Go CLI built on top of `github.com/orbcorp/orb-go`.

## Layers

```text
cmd/orb            process entrypoint
internal/cli       argv parsing and validation
internal/app       command behavior over small interfaces
internal/orbclient production adapter around orb-go
internal/output    stable JSON envelope
```

The app layer depends on interfaces defined at the consumer side. Tests inject fakes and do not require network access or an Orb API key.

## Output

All commands return:

```json
{
  "success": true,
  "data": {},
  "meta": {},
  "error": null
}
```

Errors use the same envelope and redact obvious credential markers.

## Testing

Use TDD for each command family:

- Parser tests for CLI shape.
- App tests with fake Orb clients.
- Output tests for JSON contract and redaction.
- Adapter tests with `httptest.Server` when SDK request serialization matters.

No end-to-end tests should hit Orb.
