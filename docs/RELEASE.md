# Release Guide

## Versioning

Use semantic version tags:

```sh
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

Pushing a `v*` tag triggers `.github/workflows/release.yml`.

## What Release Automation Does

1. Runs `go test ./...`.
2. Builds `orb` for:
   - Linux `amd64`, `arm64`
   - macOS `amd64`, `arm64`
   - Windows `amd64`
3. Packages archives:
   - `orb-linux-amd64.tar.gz`
   - `orb-linux-arm64.tar.gz`
   - `orb-darwin-amd64.tar.gz`
   - `orb-darwin-arm64.tar.gz`
   - `orb-windows-amd64.zip`
4. Generates `checksums.txt`.
5. Signs archives and checksums with Sigstore cosign through GitHub OIDC.
6. Publishes a GitHub Release.
7. Updates `vicentereig/homebrew-tap` formula `orb.rb`.

## Homebrew

After release:

```sh
brew install vicentereig/tap/orb
```

The release workflow expects `HOMEBREW_TAP_TOKEN` to have push access to `vicentereig/homebrew-tap`.

## Go Install

```sh
go install github.com/vicentereig/orb/cmd/orb@v0.1.0
```
