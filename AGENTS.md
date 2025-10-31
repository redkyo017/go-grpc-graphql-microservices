# Repository Guidelines

## Project Structure & Module Organization
The repo hosts four Go services: `account`, `catalog`, and `order` live under matching directories with `cmd/<service>/main.go` entrypoints, gRPC handlers in `server.go`, and repositories plus generated stubs in `pb/`. Database seed migrations live beside each service (`up.sql`, `db.dockerfile`). The `graphql` gateway wraps these services via gqlgen (`schema.graphql`, `generated.go`). Shared tooling and vendored modules sit in `pkg/`, while top-level `docker-compose.yaml` orchestrates the stack.

## Build, Test, and Development Commands
- `go build ./...` validates all modules compile with the current toolchain.
- `go test ./...` should pass before every push; add targeted package paths if runs become slow.
- `go run ./graphql` starts the GraphQL gateway on :8080 outside Docker; use `ACCOUNT_SERVICE_URL` etc. to point at running services.
- `docker-compose up --build` launches the full microservice suite plus Postgres and Elasticsearch backends.
- Regenerate protobufs with `go generate ./account/...` and `go generate ./catalog/...` before committing proto changes.

## Coding Style & Naming Conventions
Format Go code with `gofmt` (run `gofmt -w` on touched files) and keep imports ordered via `goimports`. Packages use lower_snake directory names matching import paths, exported structs follow PascalCase, and gRPC/GraphQL handlers stick to the `FooService` naming used in the existing stubs. Keep configuration structs tagged with `envconfig` for clarity.

## Testing Guidelines
Add `_test.go` files alongside the package they cover and favor table-driven tests. When introducing new behaviors, include integration coverage that can run against `docker-compose` using ephemeral containers, but gate long-running suites behind build tags. Document any schema changes with updated snapshots or fixtures.

## Commit & Pull Request Guidelines
Use short, imperative commit subjects prefixed with the touched area (example: `account: add balance repository`). PRs should describe the service impact, list config or schema changes, and attach evidence of `go test ./...` plus relevant manual checks (curl, GraphQL queries, screenshots). Link issues when available and call out follow-up work or migrations.
