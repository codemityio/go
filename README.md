# ![Go](logo.jpg)

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-91%25-brightgreen.svg?longCache=true&style=flat)

## Table of contents

- [Summary](#summary)
- [Development](#development)
- [Packages](#packages)
- [License](#license)

## Summary

Go shared library.

## Development

To work with the codebase, use `make` command as the primary entry point for all project tools.

Use the arrow keys `↓ ↑ → ←` to navigate the options, and press `/` to toggle search.

## Packages

- [`dsn`](pkg/dsn/README.md) - A lightweight, scheme-aware parser that normalizes **DSN** strings from **URLs** and file
  paths into a structured form.
- [`env`](pkg/env/README.md) - This package provides a generic, type-safe wrapper around
  [`go-envconfig`](https://github.com/sethvargo/go-envconfig) for loading environment variables into strongly-typed Go
  structs. It exposes two factory functions: `NewEnv[T]()` for error-handled config loading, and `MustBuildEnv[T]()` for
  startup-time initialization that panics on misconfiguration.
- [`fsm`](pkg/fsm/README.md) - This package contains a simple **Finite State Machine** package.
- [`json`](pkg/json/README.md) - This package provides generic, type-safe utilities for unmarshalling raw JSON into
  typed Go structs, along with testing helpers for managing deterministic sequences of JSON values.
- [`logger`](pkg/logger/README.md) - This package provides a simple asynchronous logger that processes log entries
  through a buffered channel and a background worker. It supports multiple log levels (`debug`, `info`, `warn`, `error`,
  `fatal`) and both `text` and `json` output formats, with structured fields for contextual metadata.
- [`strings`](pkg/strings/README.md) - This package provides string manipulation utilities.
- [`sync`](pkg/sync/README.md) - This package provides concurrency primitives and thread-safe data structures to
  complement the Go standard library.
- [`time`](pkg/time/README.md) - This package provides time utilities, currently offering a replaceable time provider to
  simplify time-dependent logic in tests.
- [`validator`](pkg/validator/README.md) - This package provides JSON schema-based validation for **JSON**, **YAML**,
  and arbitrary Go values.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
