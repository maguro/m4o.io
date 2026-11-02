# AGENTS.md

## Required reads

Before making any code or test edits, read `CODING_STANDARDS.md`. It holds the
rules that apply to every edit. If this file and `CODING_STANDARDS.md` differ,
follow the stricter rule.

## Skills

Use these skills from `.claude/skills/` when the task matches:

| Skill | Use it when |
|---|---|
| `go-testing` | Writing or changing a test |
| `git-commit-message` | Writing a commit message |
| `asd-ste100` | Writing or checking a code comment or a commit body |

## Project overview

- Language: Go
- Module: `m4o.io/vanity`, a `package main` binary, not a library
- What the code does, and how to deploy it: see `README.md`
- Go version: see `go.mod`. The builder image in the `Dockerfile` must be the
  same version.

`main.go` holds the whole program. It reads the import paths from a TOML file
and serves them with the handler of the `l7e.io/vanity` library. The library
holds the request handling.

## Commands

This repo has no Makefile. Use the Go tools directly.

- Format: `gofmt -w .`
- Format check: `gofmt -l .`, which prints nothing when the tree is clean
- Build: `go build ./...`
- Vet: `go vet ./...`
- Unit tests, with the race detector: `go test -race ./...`
- Lint: `golangci-lint run ./...`
- Vulnerability scan: `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`
- Tidy modules: `go mod tidy`
- Tidy check: `go mod tidy -diff`

There is no `.golangci.yml`, so `golangci-lint` runs with its default linters.

### Gate

A change that touches code is not done until each of these passes in the
working tree:

```sh
gofmt -l .
go build ./...
go vet ./...
go test -race ./...
golangci-lint run ./...
go mod tidy -diff
```

## Compatibility surface

The binary exports no symbol, so the public API rule in `CODING_STANDARDS.md`
has nothing to apply to. These are the contracts that a change must not break
without the same care:

- The environment variables `PORT`, `TOML_FILE`, `TOML_TABLE`, and `LOG_LEVEL`,
  and their defaults. The `Dockerfile` sets the first three, and the Cloud Run
  service can set any of them.
- The schema of `config.toml`: a table of entries, each with `import_path`,
  `vcs`, and `vcs_path`.
- The `go-import` meta tag of the response. The Go toolchain reads it.

## Continuous integration

This repo has no GitHub Actions workflow. `.travis.yml` is present, but it
builds Go 1.11 to 1.14, which `go.mod` no longer permits. Treat the gate above
as the only check, and run it by hand.
