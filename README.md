# Verified ID Terraform Provider

This is a fork from https://github.com/microsoft/terraform-provider-msgraph

Alpha version of a Terraform provider for Microsoft Entra Verified ID.

## DEVELOPMENT

### Prerequisites

- [Go](https://go.dev/dl/) 1.22+ (check [go.mod](go.mod) for the exact version)
- [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.5+
- (Optional) `make` — on Windows install via [Chocolatey](https://chocolatey.org/) (`choco install make`), or use [Git Bash](https://git-scm.com/downloads) / WSL

> Note for C# developers: the `go` CLI is roughly the equivalent of `dotnet` — it builds, tests, and manages dependencies. The committed `vendor/` folder is similar to a checked-in `packages/` folder.

### Build

Build the provider binary:

```sh
go build -o terraform-provider-verifiedid
```

Install the locally built provider into the Terraform plugin folder so the configs in [examples/](examples/) can pick it up:

**Linux / macOS**
```sh
cp terraform-provider-verifiedid ~/.terraform.d/plugins/local/custom/verifiedid/1.0.0/linux_amd64/
```

**Windows (PowerShell)**
```powershell
$dest = "$env:APPDATA\terraform.d\plugins\local\custom\verifiedid\1.0.0\windows_amd64"
New-Item -ItemType Directory -Force -Path $dest | Out-Null
Copy-Item .\terraform-provider-verifiedid.exe $dest -Force
```

### Run tests

Go tests work like xUnit/NUnit: any `func TestXxx(t *testing.T)` in a `_test.go` file is a test. This repo has two flavors:

- **Unit tests** — fast, no network, no Azure tenant required.
- **Acceptance tests** — real HTTP calls to Microsoft Graph; require auth (`ARM_*` env vars) and `TF_ACC=1`. They take time and may create/delete real resources in your tenant.

#### Unit tests (recommended starting point)

Run all unit tests in the module:

```sh
go test ./... -v
```

Run a single package (like running tests in one C# project):

```sh
go test ./internal/services/... -v
```

Run a single test by name (regex match, similar to `dotnet test --filter`):

```sh
go test ./internal/services -run TestAccVerifiedIDResource_basic -v
```

Via the Makefile (Linux/macOS/WSL/Git Bash):

```sh
make test
```

#### Acceptance tests

Configure credentials first (see [docs/index.md](docs/index.md) for all auth options), then:

**Linux / macOS / WSL**
```sh
export ARM_TENANT_ID=...
export ARM_CLIENT_ID=...
export ARM_CLIENT_SECRET=...
export TF_ACC=1
make testacc
```

**Windows (PowerShell)**
```powershell
$env:ARM_TENANT_ID = "..."
$env:ARM_CLIENT_ID = "..."
$env:ARM_CLIENT_SECRET = "..."
$env:TF_ACC = "1"
go test ./internal/services/... -v -timeout 300m
```

Run a single acceptance test:

```sh
go test ./internal/services -run TestAccVerifiedIDResource_basic -v -timeout 60m
```

### Lint (golangci-lint)

This repo uses [`golangci-lint`](https://golangci-lint.run/) (config: [.golangci.yml](.golangci.yml)) — analogous to Roslyn analyzers in C#.

Install it once:

**Linux / macOS / WSL**
```sh
# pinned binary install (matches the version used by `make tools`)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin v1.61.0
# or via Homebrew
brew install golangci-lint
```

**Windows (PowerShell)**
```powershell
# via Scoop
scoop install golangci-lint
# or via Chocolatey
choco install golangci-lint
# or with go install (matches the pinned version)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
```

Make sure `$(go env GOPATH)/bin` (Linux/macOS) or `%USERPROFILE%\go\bin` (Windows) is on your `PATH`.

Run the linter against the whole module:

```sh
golangci-lint run ./...
```

Or via the Makefile (Linux/macOS/WSL/Git Bash):

```sh
make lint
```

Lint a single package while iterating:

```sh
golangci-lint run ./internal/services/...
```

Auto-fix the issues that support it:

```sh
golangci-lint run ./... --fix
```

### Other useful developer commands

| Task | Command | C# analogue |
|---|---|---|
| Format code | `go fmt ./...` or `make fmt` | `dotnet format` |
| Stricter format | `gofumpt -w .` | — |
| Lint | `make lint` (uses `golangci-lint`) | Roslyn analyzers |
| Tidy dependencies | `go mod tidy` | `dotnet restore` cleanup |
| Refresh `vendor/` | `go mod vendor` | restore packages |
| Verify go.mod / vendor are clean | `make depscheck` | — |
| Install dev tooling (linters, terrafmt, …) | `make tools` | one-time tool install |
| Regenerate docs | `make docs` (runs `go generate`) | — |
| Format `.tf` examples & docs | `make terrafmt` | — |
| Lint `.tf` examples | `make tflint` | — |
| Run all pre-commit checks | `make fmtcheck` | — |

### Logging

Enable verbose Terraform/HTTP logging when running `terraform` against the locally built provider:

**Linux / macOS / WSL**
```sh
export TF_LOG=DEBUG   # or WARN, INFO, TRACE
```

**Windows (PowerShell)**
```powershell
$env:TF_LOG = "DEBUG"
```

### Project layout (quick map)

- [main.go](main.go) — entry point (registers the provider, like `Program.cs`).
- [internal/provider/](internal/provider/) — provider wiring and authentication.
- [internal/services/](internal/services/) — resources & data sources (the meat of the provider).
- [internal/clients/](internal/clients/) — Microsoft Graph HTTP client.
- [docs/](docs/), [templates/](templates/) — generated and templated user documentation.
- [examples/](examples/) — sample Terraform configurations used in docs and manual testing.
- [vendor/](vendor/) — checked-in third-party Go dependencies.