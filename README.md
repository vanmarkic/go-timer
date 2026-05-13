# go-timer

A small, portable countdown timer with second precision. Single Go binary that serves a local web UI with **Start**, **Stop**, **Reset**, and **Set Time** buttons.

## Download

Prebuilt binaries for Linux, macOS, and Windows (amd64 and arm64) are published on the [Releases page](../../releases). Each release includes:

- `go-timer_<version>_<os>_<arch>.tar.gz` (or `.zip` on Windows)
- `checksums.txt` — SHA-256 of every asset
- `*.spdx.json` — SPDX SBOM per archive
- Build provenance attestations (verifiable with `gh attestation verify`)

## Run

```sh
./go-timer            # listen on 127.0.0.1:8080 and open a browser
./go-timer -addr :9000 -no-open
```

The Go binary embeds the UI (`web/*`) and serves it locally; all timer logic runs in the browser, so there is no API and no network traffic leaves the loopback interface.

Keyboard shortcuts: **Space** = start/stop, **R** = reset.

## Build from source

```sh
go build -trimpath -ldflags="-s -w" -o go-timer .
```

`CGO_ENABLED=0` is the default for releases, producing a static binary.

## Supply-chain security

Every push and every release tag runs the following gates — a release will not publish if any of them fail at HIGH+ severity:

| Tool | What it checks |
|------|----------------|
| `go mod verify` | Module checksums match `go.sum` |
| [`govulncheck`](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) | Known vulnerabilities in stdlib and deps actually reachable from your code |
| [`gosec`](https://github.com/securego/gosec) | Go-specific SAST rules |
| [`trivy`](https://github.com/aquasecurity/trivy) | Filesystem vulnerability + secret + misconfig scan |
| [`syft`](https://github.com/anchore/syft) | Generates SPDX + CycloneDX SBOMs |
| [`grype`](https://github.com/anchore/grype) | Vulnerability scan of the generated SBOM |
| `actions/attest-build-provenance` | SLSA build provenance attestation on release artifacts |
| Dependabot | Weekly updates for `gomod` and `github-actions` |

Workflows live in [`.github/workflows/`](.github/workflows/).

## Release

Push a tag of the form `vX.Y.Z`:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The `release` workflow runs the pre-flight security gate, then GoReleaser cross-compiles, generates per-archive SBOMs and SHA-256 checksums, and uploads everything to a GitHub Release.
