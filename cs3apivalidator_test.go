package main

import (
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/scenario"
	flag "github.com/spf13/pflag"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	cfg := featurecontext.Config{}
	flag.StringVar(&cfg.Endpoint, "endpoint", "localhost:9142", "Endpoint Url and port of a running cs3 implementation")
	flag.StringVar(&cfg.GrpcTLSMode, "grpc-tls-mode", "off", "TLS mode for grpc client connections ('off', 'on' or 'insecure')")
	flag.BoolVar(&cfg.AsyncPropagation, "async-propagation", false, "Enable async propagation")
	flag.DurationVar(&cfg.AsyncPropagationDelay, "async-propagation-delay", 200*time.Millisecond, "Delay for async propagation")
	flag.BoolVar(&cfg.HttpInsecure, "http-insecure", true, "Allow insecure HTTP connections")
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "cs3apiValidator",
		ScenarioInitializer: scenario.InitializeScenario(cfg),
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
