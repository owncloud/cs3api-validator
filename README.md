# cs3api-validator

<img width="100px" src="https://raw.githubusercontent.com/cs3org/logos/master/cs3org/cs3org.png"/>

[![Build Status](https://drone.owncloud.com/api/badges/owncloud/cs3api-validator/status.svg)](https://drone.owncloud.com/owncloud/cs3api-validator)
[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## End-to-End Test Suite for the CS3 APIs

**This tool will receive a lot of changes before version 1.0.0**

The cs3api-validator is a tool to test implementations of the [CS3Apis](https://github.com/cs3org/cs3apis). It works as standalone software which only needs the address of a running cs3api provider.

## Purpose

### BDD (Behavior driven development)

The cs3api-validator can be run locally in a development phase to develop against a well-defined set of basic API operations. It has no external dependencies and runs human-readable gherkin test scenarios. This helps to understand the behavior of the CS3APIs being an additional way of documenting the API in combination with the specification.

### Litmus Testing

This tool makes it possible to confirm that an implementation of the CS3APIs is compliant to the spec and fulfills the basic operations. This helps the CS3 community to keep different implementations in sync and foster compatibility between them.

## Contributions

This is a community driven Open Source Project. We welcome contributions from everyone and we're ready to support you if you have the enthusiasm to contribute.

Please use the [issue tracker](https://github.com/owncloud/cs3api-validator/issues) to report problems or propose changes.

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

You can run the tests with the built-in go test command. The command passes its flags to the godog test suite. The test suite needs one flag to be set: The network address of the running system under test. It defaults to `localhost:9142` and you can set it using the ``--endpoint``flag.

> **_NOTE:_** If you want to use the godog flags you need to prefix them with ``godog.flagname``.

#### In a working go environment

Run ``go test --endpoint=your-addr:port -v``

#### Build a binary with the tests

Run ``go test -c``. This will create a ``cs3api-validator.test`` binary.

Execute the tests ``./cs3qpi-validator.test --endpoint=your-addr:port``

### Use tags

You can use [tags](https://github.com/cucumber/godog#tags) to filter features which should be executed. `--godog.tags=<expression>`

## License

Apache-2.0


