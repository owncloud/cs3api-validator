package scenario

import (
	"github.com/cucumber/godog"
	fc "github.com/owncloud/cs3api-validator/featurecontext"
)

// Endpoint the address of the grpc cs3api provider
var config fc.Config

// InitializeScenario wraps to scenario initialization to pass configuration
func InitializeScenario(c fc.Config) func(*godog.ScenarioContext) {
	config = c
	return initializeScenario
}

func initializeScenario(sc *godog.ScenarioContext) {
	f := newFeatureContext(sc)
	f.Init(config)
}
