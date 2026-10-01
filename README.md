# E2Engine CLI

Command-line interface for E2Engine.

This module provides the E2Engine command-line application for creating, validating, listing, inspecting, deleting, and running E2Engine resources using local E2Engine services.

## Build

```bash
make build-cli
```

## Packages

**cmd/cli**

E2Engine CLI entry point and Cobra command definitions.

**internal/command**

Command handlers for E2Engine resources and operations.

**internal/service**

CLI service construction and integration with E2Engine Core and persistence implementations.

**internal/output**

Output configuration and rendering used by CLI commands.

**internal/runnerprocess**

Local runner process lifecycle and launching.

## Development

Run tests:

```bash
make test
```

Run the linter:

```bash
make lint
```

Run race detection:

```bash
make test-race
```

Run the complete verification suite:

```bash
make verify
```

## E2Engine

This repository is part of E2Engine.

- core — core domain model, execution logic, and public APIs
- repository — persistence implementations
- runner-local — local test execution
- cli — command-line interface
- tests — end-to-end tests for E2Engine
- demo — executable demonstration system and E2Engine examples

## License

Licensed under the Apache License, Version 2.0.