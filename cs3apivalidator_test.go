package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/owncloud/cs3api-validator/scenario"
	flag "github.com/spf13/pflag"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

// Endpoint GRPC address of a running CS3 implementation
var Endpoint string

// HTTPinsecure controls wether insecure HTTP connections are allowed or not
var HTTPinsecure bool

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.StringVar(&Endpoint, "endpoint", "localhost:9142", "Endpoint Url and port of a running cs3 implementation")
	flag.BoolVar(&HTTPinsecure, "http-insecure", true, "Allow insecure HTTP connections")
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "cs3apiValidator",
		ScenarioInitializer: scenario.InitializeScenario(Endpoint, HTTPinsecure),
		Options:             &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
