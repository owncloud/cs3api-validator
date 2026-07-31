# CS3 API Validator

<!-- OSPO-managed README | Generated: 2026-04-16 | v2 -->

[![License](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](LICENSE) [![ownCloud OSPO](https://img.shields.io/badge/OSPO-ownCloud-blue)](https://kiteworks.com/opensource) [![Docker Hub](https://img.shields.io/docker/pulls/owncloud)](https://hub.docker.com/r/owncloud/ocis)

The CS3 API Validator is an end-to-end test suite for implementations of the [CS3 APIs](https://github.com/cs3org/cs3apis). It runs human-readable Gherkin test scenarios against a CS3 API provider, serving both as a BDD development tool and as a litmus test to verify that an implementation complies with the CS3 API specification. The tool requires only the network address of a running CS3 API provider and has no external dependencies beyond Go.

## Part of oCIS

This tool is part of the [ownCloud Infinite Scale (oCIS)](https://github.com/owncloud/ocis) ecosystem and is used to validate CS3 API compliance. It helps keep different CS3 API implementations in sync and fosters interoperability across the [CS3 community](https://cs3community.org/).

This component is part of the [oCIS Docker image](https://hub.docker.com/r/owncloud/ocis).

## Getting Started

Follow the steps below to run the CS3 API validation suite.

### Quick Start

```bash
git clone git@github.com:owncloud/cs3api-validator.git
cd cs3api-validator
go test -v  # default network address is localhost:9142
```

### Adding New Test Features

Add Gherkin feature files to the `features/` directory, then implement the step definitions in Go. Run `go run github.com/cucumber/godog/cmd/godog@master` to see which steps need implementation.

## Usage

Instructions for running and configuring the CS3 API validation suite:

### Running Tests

Run with the built-in `go test` command. The `--endpoint` flag sets the network address of the system under test (defaults to `localhost:9142`):

```bash
go test --endpoint=your-addr:port -v
```

To build a standalone test binary:

```bash
go test -c
./cs3api-validator.test --endpoint=your-addr:port
```

### Filtering with Tags

Use [Godog tags](https://github.com/cucumber/godog#tags) to select which features to run:

```bash
go test --godog.tags="@smoke" -v
```

### Purpose

- **BDD development**: Run locally during development against a well-defined set of basic CS3 API operations. Human-readable Gherkin scenarios serve as both tests and API documentation.
- **Litmus testing**: Confirm that a CS3 API implementation is spec-compliant and supports basic operations, helping keep different implementations in sync.

### Adding New Features

1. Add Gherkin feature files to the `features/` directory
2. Run `go run github.com/cucumber/godog/cmd/godog@master` to get step definition stubs
3. Implement steps in a `*_test.go` file using the `FeatureContext` struct as a receiver for sharing state between steps

## Documentation

- [CS3 API specification](https://github.com/cs3org/cs3apis)
- [Godog (Go Cucumber)](https://github.com/cucumber/godog)
- [Gherkin syntax](https://cucumber.io/docs/gherkin/)

## Community & Support

**[Star](https://github.com/owncloud/cs3api-validator)** this repo and **Watch** for release notifications!

- [ownCloud Website](https://owncloud.com)
- [Community Discussions](https://github.com/orgs/owncloud/discussions)
- [Matrix Chat](https://app.element.io/#/room/#owncloud:matrix.org)
- [Documentation](https://doc.owncloud.com)
- [Enterprise Support](https://owncloud.com/contact-us/)
- [OSPO Home](https://kiteworks.com/opensource)

## Contributing

We welcome contributions! Please read the [Contributing Guidelines](CONTRIBUTING.md)
and our [Code of Conduct](CODE_OF_CONDUCT.md) before getting started.

### Workflow

- **Rebase Early, Rebase Often!** We use a rebase workflow. Always rebase on the target branch before submitting a PR.
- **Dependabot**: Automated dependency updates are managed via Dependabot. Review and merge dependency PRs promptly.
- **Signed Commits**: All commits **must** be PGP/GPG signed. See [GitHub's signing guide](https://docs.github.com/en/authentication/managing-commit-signature-verification).
- **DCO Sign-off**: Every commit must carry a `Signed-off-by` line:
  ```
  git commit -s -S -m "your commit message"
  ```
- **GitHub Actions Policy**: Workflows may only use actions that are (a) owned by `owncloud`, (b) created by GitHub (`actions/*`), or (c) verified in the GitHub Marketplace.

## Security

**Do not open a public GitHub issue for security vulnerabilities.**

Report vulnerabilities at **<https://security.owncloud.com>** -- see [SECURITY.md](SECURITY.md).

Bug bounty: [YesWeHack ownCloud Program](https://yeswehack.com/programs/owncloud-bug-bounty-program)

## License

This project is licensed under the [Apache-2.0](LICENSE).

## About the ownCloud OSPO

The [Kiteworks Open Source Program Office](https://kiteworks.com/opensource), operating under
the [ownCloud](https://owncloud.com) brand, launched on May 5, 2026, to steward the open source
ecosystem around ownCloud's products. The OSPO ensures transparent governance, license compliance,
community health, and sustainable collaboration between the open source community and
[Kiteworks](https://www.kiteworks.com), which acquired ownCloud in 2023.

- **OSPO Home**: <https://kiteworks.com/opensource>
- **GitHub**: <https://github.com/owncloud>
- **ownCloud**: <https://owncloud.com>

For questions about the OSPO or licensing, contact ospo@kiteworks.com.

> **License status:** This repository is already licensed under Apache-2.0 -- the OSPO target license.
> No migration is required.
