# AI Agent Guidelines for CS3 API Validator

This file provides context for AI coding agents (Claude Code, GitHub Copilot, Cursor, etc.) working in this repository.

## Repository Overview
- **Product family:** oCIS
- **Primary language(s):** Go, Gherkin
- **Build system:** Go modules
- **Test framework:** Godog (Go Cucumber), Go testing
- **CI system:** None detected

## Architecture & Key Paths
- `features/` - Gherkin feature files (BDD test scenarios)
- `featurecontext/` - Feature context setup
- `steps/` - Go step definitions for Gherkin scenarios
- `helpers/` - Test helper utilities
- `constants/` - Shared constants
- `scenario/` - Scenario management
- `cs3apivalidator.go` - Main validator logic
- `cs3apivalidator_test.go` - Test entry point
- `go.mod` - Go module definition
- `go.sum` - Dependency checksums
- `docker/` - Docker configuration

## Development Conventions
- **Branching:** main
- **Commit messages:** DCO sign-off required (`git commit -s`)
- **Code style:** No specific linter configured
- **PR process:** Open a PR against main. All CI checks must pass.

## Build & Test Commands
```bash
# Build
go build ./...

# Test (requires running CS3 API provider at localhost:9142)
go test -v

# Run with custom address
go test -v -addr=myserver:9142

# Lint
Not detected
```

## Important Constraints
- All code contributions must be compatible with the **Apache-2.0** license
- Do not introduce new **copyleft-licensed dependencies** (GPL, AGPL, LGPL, MPL) without explicit discussion in an issue first. This is especially important for repos migrating to Apache 2.0.
- Do not introduce new dependencies without discussion in an issue first
- Tests require a running CS3 API provider (e.g., oCIS or Reva)


## OSPO Policy Constraints

### GitHub Actions
- **Only** use actions owned by `owncloud`, created by GitHub (`actions/*`), verified on the GitHub Marketplace, or verified by the ownCloud Maintainers.
- Pin all actions to their full commit SHA (not tags): `uses: actions/checkout@<SHA> # vX.Y.Z`
- Never introduce actions from unverified third parties.

### Dependency Management
- Dependabot is configured for automated dependency updates.
- Review and merge Dependabot PRs as part of regular maintenance.
- Do not introduce new dependencies without discussion in an issue first.

### Git Workflow
- **Rebase policy**: Always rebase; never create merge commits. Use `git pull --rebase` and `git rebase` before pushing.
- **Signed commits**: All commits **must** be PGP/GPG signed (`git commit -S -s`).
- **DCO sign-off**: Every commit needs a `Signed-off-by` line (`git commit -s`).
- **Conventional Commits & Squash Merge**: Use the [Conventional Commits](https://www.conventionalcommits.org/) format where the repository enforces it. Many repos use squash merge, where the PR title becomes the commit message on the default branch — apply Conventional Commits format to PR titles as well. A reusable GitHub Actions workflow enforces this.

## Context for AI Agents
- Match existing code style
- Do not refactor unrelated code in the same PR
- Write tests for new functionality
- Keep PRs focused and atomic
- New test scenarios should be added as Gherkin feature files in `features/`
- Step definitions are written in Go - follow existing patterns
