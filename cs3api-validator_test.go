package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts) // godog v0.11.0 and later
}

func userHasCreatedAPersonalSpace(user string) error {
	return godog.ErrPending
}

func userListsAllAvailableSpaces(user string) error {
	return godog.ErrPending
}

func onePersonalSpaceShuoldBeListedInTheResponse() error {
	return godog.ErrPending
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^user "([^"]*)" has created a personal space$`, userHasCreatedAPersonalSpace)
	ctx.Step(`^user "([^"]*)" lists all available spaces$`, userListsAllAvailableSpaces)
	ctx.Step(`^one personal space should be listed in the response$`, onePersonalSpaceShuoldBeListedInTheResponse)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                 "cs3api-validator",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
