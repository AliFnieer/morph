# Contributing

Thanks for your interest in improving morph.

## Workflow

1. Create a feature branch from `dev`.
2. Make your changes and add tests where appropriate.
3. Open a pull request against `dev`.
4. Once `dev` is stable, it is merged into `main`.

`main` is the release branch and should always be green.

## Development

```bash
go build ./...
go test ./...
go vet ./...
gofmt -l .
```

Integration tests build the `morph` binary automatically via `tests/main_test.go`.

## Commit messages

Use clear, imperative messages, for example:

- `feat: add webp support`
- `fix: handle empty input files`
- `test: cover archive conversions`
- `docs: update install instructions`

## Releases

Releases are cut from `main` by pushing a version tag (for example `v0.1.0`).
The `Release` GitHub Actions workflow then builds cross-platform binaries with
GoReleaser and publishes them to GitHub Releases.
