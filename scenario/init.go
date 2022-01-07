package scenario

import (
	"github.com/cucumber/godog"
)

// available settings
var Endpoint string
var HttpInsecure bool

// InitializeScenario wraps to scenario initialization to pass configuration
func InitializeScenario(endpoint string, httpInsecure bool) func(*godog.ScenarioContext) {
	Endpoint = endpoint
	HttpInsecure = httpInsecure
	return initializeScenario
}

func initializeScenario(sc *godog.ScenarioContext) {
	f := newFeatureContext()
	f.Init(
		Endpoint,
		HttpInsecure,
	)
	f.registerSteps(sc)
}
