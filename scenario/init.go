package scenario

import (
	"github.com/cucumber/godog"
)

// Endpoint the address of the grpc cs3api provider
var Endpoint string
// HttpInsecure flag to use insecure transport
var HttpInsecure bool

// InitializeScenario wraps to scenario initialization to pass configuration
func InitializeScenario(endpoint string, httpInsecure bool) func(*godog.ScenarioContext) {
	Endpoint = endpoint
	HttpInsecure = httpInsecure

	return initializeScenario
}

func initializeScenario(sc *godog.ScenarioContext) {
	f := newFeatureContext(sc)
	f.Init(
		Endpoint,
		HttpInsecure,
	)
}
