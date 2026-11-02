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

- Build: `make build`
- Unit tests, with the race detector: `make test`
- Coverage: `make coverage`
- Lint: `make lint`
- Format: `make fmt`
- Format check: `make check-fmt`
- Vulnerability scan: `make vuln`
- Tidy modules: `make tidy`
- Gate target: `make check`
- CI target, the gate and coverage: `make ci`

`make check` runs `test`, `lint`, and `check-fmt`. The lint and format targets
install the pinned `golangci-lint` into `.cache/`, which `.gitignore` excludes.
`.golangci.yml` holds the linter set and the import grouping of this module.

### Gate

A change that touches code is not done until `make check` passes in the
working tree.

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
builds Go 1.11 to 1.14, which `go.mod` no longer permits. Run `make check` by
hand until a workflow replaces it.
