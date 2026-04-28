# orb Summary

`orb` is a Go CLI for Orb billing forensics.

Current implemented surface:

- `orb version`
- `orb ping`

Design goals:

- Read-heavy investigation workflows.
- JSON-first output for humans, scripts, Codex, and Claude.
- No real Orb API keys in tests.
- Release as `orb` through GitHub Releases and `vicentereig/homebrew-tap`.
