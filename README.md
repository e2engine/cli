# E2Engine CLI

Command-line interface for E2Engine.

Part of [E2Engine](https://e2engine.dev), an open-source platform for declarative end-to-end testing.

Detailed documentation for the E2Engine CLI is available at [https://e2engine.dev/docs/using-e2engine/cli/](https://e2engine.dev/docs/using-e2engine/cli/).

This module provides the E2Engine command-line application for managing and running E2Engine resources using local E2Engine services.

## Installation

### macOS

Install E2Engine using Homebrew:

```bash
brew install --cask e2engine/tap/e2engine
```

### Linux

Download the appropriate `.deb` or `.rpm` package for your architecture from the [latest GitHub release](https://github.com/e2engine/cli/releases/latest).

Debian/Ubuntu:

```bash
sudo apt install ./e2engine_*.deb
```

Fedora/RHEL:

```bash
sudo dnf install ./e2engine_*.rpm
```

### Windows

Install E2Engine using Scoop:

```powershell
scoop bucket add e2engine https://github.com/e2engine/scoop-bucket
scoop install e2engine
```

## Build

```bash
make build-cli
```

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

- [core](https://github.com/e2engine/core) — core domain model, execution logic, and public APIs
- [repository](https://github.com/e2engine/repository) — persistence implementations
- [runner-local](https://github.com/e2engine/runner-local) — local test execution
- [cli](https://github.com/e2engine/cli) — command-line interface
- [tests](https://github.com/e2engine/tests) — end-to-end tests for E2Engine
- [demo](https://github.com/e2engine/demo) — executable demonstration system and E2Engine usage examples
- [instrumentation-go](https://github.com/e2engine/instrumentation-go) — Go instrumentation library for E2Engine

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).