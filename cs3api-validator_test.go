package main

import (
	"os"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
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

type User struct {
	RevaToken string
}

type FeatureContext struct {
	Client gateway.GatewayAPIClient
	Users  []User
}

func (f *FeatureContext) userHasCreatedAPersonalSpace(user string) error {
	return godog.ErrPending
}

func (f *FeatureContext) userListsAllAvailableSpaces(user string) error {
	return godog.ErrPending
}

func (f *FeatureContext) onePersonalSpaceShuoldBeListedInTheResponse() error {
	return godog.ErrPending
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	f := &FeatureContext{}
	ctx.Step(`^user "([^"]*)" has created a personal space$`, f.userHasCreatedAPersonalSpace)
	ctx.Step(`^user "([^"]*)" lists all available spaces$`, f.userListsAllAvailableSpaces)
	ctx.Step(`^one personal space should be listed in the response$`, f.onePersonalSpaceShuoldBeListedInTheResponse)
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
