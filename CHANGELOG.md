# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2026-11-04

### Added

- Unit tests for the registry, format detection, and text/data converters.
- Dependabot configuration for Go modules and GitHub Actions.
- Issue and pull request templates.

### Changed

- Bump GitHub Actions: `checkout` to v7, `setup-go` to v7, and
  `goreleaser-action` to v7.

### Fixed

- Resolve `golangci-lint` findings and run the linter in CI.
- Run the 13 integration tests in `tests/more_test2.go` that were previously
  skipped because the file was not named `*_test.go`.
- Report the actual version in `morph doctor` instead of a hardcoded value.

## [0.1.0] - 2026-11-04

### Added

- Offline, privacy-focused file conversion CLI with a registry-based design.
- Converters for images, text/data (json, yaml, toml, csv, ndjson), html/markdown,
  pdf, office (docx, xlsx, pptx, odt, ods, odp), epub, rtf, svg, and archives
  (zip, tar, gz, tar.xz, tar.bz2).
- Batch/glob and recursive directory processing.
- Stdin/stdout support for piping data.
- `doctor`, `version`, `completion`, `--list`, and `--list-json` commands.
- Verbose/quiet/force flags, PDF text extraction, and fuzzy conversion suggestions.
- CI workflow, GoReleaser release configuration, and integration test suite.

[Unreleased]: https://github.com/AliFnieer/morph/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/AliFnieer/morph/releases/tag/v0.1.1
[0.1.0]: https://github.com/AliFnieer/morph/releases/tag/v0.1.0
