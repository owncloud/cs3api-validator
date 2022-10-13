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

// endpoint GRPC address of a running CS3 implementation
var endpoint string

// httpInsecure controls whether insecure HTTP connections are allowed or not
var httpInsecure bool

// grpcTLSMode TLS mode for grpc client connections
var grpcTLSMode string

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.StringVar(&endpoint, "endpoint", "localhost:9142", "Endpoint Url and port of a running cs3 implementation")
	flag.StringVar(&grpcTLSMode, "grpc-tls-mode", "off", "TLS mode for grpc client connections ('off', 'on' or 'insecure')")
	flag.BoolVar(&httpInsecure, "http-insecure", true, "Allow insecure HTTP connections")
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "cs3apiValidator",
		ScenarioInitializer: scenario.InitializeScenario(endpoint, httpInsecure, grpcTLSMode),
		Options:             &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

// Empty test to suppress the "no tests found" warning
func TestOne(t *testing.T) {
}
