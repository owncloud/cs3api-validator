package scenario

import (
	"github.com/cucumber/godog"
)

// Endpoint the address of the grpc cs3api provider
var endpoint string

// HttpInsecure flag to use insecure transport
var httpInsecure bool

// TLS mode for grpc client connections
var grpcTLSMode string

// InitializeScenario wraps to scenario initialization to pass configuration
func InitializeScenario(e string, h bool, g string) func(*godog.ScenarioContext) {
	endpoint = e
	httpInsecure = h
	grpcTLSMode = g

	return initializeScenario
}

func initializeScenario(sc *godog.ScenarioContext) {
	f := newFeatureContext(sc)
	f.Init(
		endpoint,
		httpInsecure,
		grpcTLSMode,
	)
}
