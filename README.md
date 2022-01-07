# cs3api-validator

<img width="100px" src="https://raw.githubusercontent.com/cs3org/logos/master/cs3org/cs3org.png"/>

[![Build Status](https://drone.owncloud.com/api/badges/owncloud/cs3api-validator/status.svg)](https://drone.owncloud.com/owncloud/cs3api-validator)
[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Rocket chat](https://img.shields.io/badge/Chat%20on%20Rocket.Chat-blue?logo=data:image/png;base64,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)](https://talk.owncloud.com/channel/ocis)

## End-to-End Test Suite for the CS3 APIs

**This tool will receive a lot of changes before version 1.0.0**

The cs3api-validator is a tool to test implementations of the [CS3Apis](https://github.com/cs3org/cs3apis). It works as a standalone software which only needs the address of a running cs3api provider.

## Purpose

### BDD (Behavior driven development)

The cs3api-validator can be run locally in a development phase to develop against a well defined set of basic API operations. It has no external dependencies and runs human readable gherkin test scenarios. This helps to understand the behavior of the CS3APIs being an additional way of documenting the API in combination with the specification.

### Litmus Testing

This tool makes it possible to confirm that an implementation of the CS3APIs is compliant to the spec and fulfills the basic operations. This helps the CS3 community to keep different implementations in sync and foster compatibility between them.

## Contributions

This is a community driven Open Source Project. We welcome contributions from everyone and we're ready to support you if you have the enthusiasm to contribute.

Please use the [issue tracker](https://github.com/owncloud/cs3api-validator/issues) to report problems or propose changes.

## Getting help

We have a [Public Chat](https://talk.owncloud.com/) where you can chat with other community members and developers.

Here are some useful channels to try:

- [#ocis](https://talk.owncloud.com/channel/ocis) - Next Generation ownCloud Architecture
- [#general](https://talk.owncloud.com/channel/general) - General Topics around File Sync & Share

## Developing

### Quick start
```shell
git clone git@github.com:owncloud/cs3api-validator.git
cd cs3api-validator
go test -v # default network addr of cs3api provider is localhost:9142
```

### Add features

Add new test steps to the feature files in the `features` directory.

```gherkin
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining
```
Then run ` go run github.com/cucumber/godog/cmd/godog@master` which will output something similar like this

```
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12          # features/godogs.feature:6
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining

1 scenarios (1 undefined)
3 steps (3 undefined)
220.129µs

You can implement step definitions for undefined steps with these snippets:

func iEat(arg1 int) error {
        return godog.ErrPending
}

func thereAreGodogs(arg1 int) error {
        return godog.ErrPending
}

func thereShouldBeRemaining(arg1 int) error {
        return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
        ctx.Step(`^I eat (\d+)$`, iEat)
        ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
        ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
```

Then copy the new step definition stubs to a *_test.go file and implement the steps. In this project we use a FeatureContext struct to share values between the tests steps. In order to do this we need to use a pointer to the FeatureContext as a function receiver in the test step methods.

```go
func (f *FeatureContext) iEat(arg1 int) error {
        return godog.ErrPending
}

func (f *FeatureContext) thereAreGodogs(arg1 int) error {
        return godog.ErrPending
}

func (f *FeatureContext) thereShouldBeRemaining(arg1 int) error {
        return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
        f := &FeatureContext{}
        ctx.Step(`^I eat (\d+)$`, f.iEat)
        ctx.Step(`^there are (\d+) godogs$`, f.thereAreGodogs)
        ctx.Step(`^there should be (\d+) remaining$`, f.thereShouldBeRemaining)
}
```

## Usage

### Run with go test

You can run the tests with the built in go test command. The command passes its flags to the godog test suite. The test suite needs one flag to be set: The network address of the running system under test. It defaults to `localhost:9142` and you can set it using the ``--endpoint``flag.

> **_NOTE:_** If you want to use the godog flags you need to prefix them with ``godog.flagname``.

#### In a working go environment

Run ``go test --endpoint=your-addr:port -v``

#### Build a binary with the tests

Run ``go test -c``. This will create a ``cs3api-validator.test`` binary.

Execute the tests ``./cs3qpi-validator.test --endpoint=your-addr:port``

### Use tags

You can use [tags](https://github.com/cucumber/godog#tags) to filter features which should be exectuted. `--godog.tags=<expression>`

## License

Apache-2.0

