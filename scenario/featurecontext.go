package scenario

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/steps/login"
	"github.com/owncloud/cs3api-validator/steps/publicshare"
	"github.com/owncloud/cs3api-validator/steps/resources"
	"github.com/owncloud/cs3api-validator/steps/spaces"
)

// featureContext embeds all available feature contexts
type featureContext struct {
	*featurecontext.FeatureContext

	*login.LoginFeatureContext
	*publicshare.PublicShareFeatureContext
	*resources.ResourcesFeatureContext
	*spaces.SpacesFeatureContext
}

// newFeatureContext returns a new feature context for the scenario initialization
// and makes sure that all contexts have the same pointer to a single FeatureContext
func newFeatureContext() *featureContext {
	fc := &featurecontext.FeatureContext{}

	// every xxxFeatureContext needs to have the pointer to a _single_ / common FeatureContext
	uc := &featureContext{
		FeatureContext: fc,

		LoginFeatureContext:       login.NewLoginFeatureContext(fc),
		PublicShareFeatureContext: publicshare.NewPublicShareFeatureContext(fc),
		ResourcesFeatureContext:   resources.NewResourcesFeatureContext(fc),
		SpacesFeatureContext:      spaces.NewSpacesFeatureContext(fc),
	}
	return uc
}

// registerSteps registers all given steps during scenario initialization
func (fc *featureContext) registerSteps(sc *godog.ScenarioContext) {
	fc.LoginFeatureContext.RegisterSteps(sc)
	fc.PublicShareFeatureContext.RegisterSteps(sc)
	fc.ResourcesFeatureContext.RegisterSteps(sc)
	fc.SpacesFeatureContext.RegisterSteps(sc)
}
