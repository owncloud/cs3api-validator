package scenario

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/steps/login"
	"github.com/owncloud/cs3api-validator/steps/publicshare"
	"github.com/owncloud/cs3api-validator/steps/resources"
	"github.com/owncloud/cs3api-validator/steps/spaces"
	"github.com/owncloud/cs3api-validator/steps/sampletest"
)

// featureContext embeds all available feature contexts
type featureContext struct {
	*featurecontext.FeatureContext

	*login.LoginFeatureContext
	*publicshare.PublicShareFeatureContext
	*resources.ResourcesFeatureContext
	*spaces.SpacesFeatureContext
	*sampletest.SampleTestFeatureContext
}

// newFeatureContext returns a new feature context for the scenario initialization
// and makes sure that all contexts have the same pointer to a single FeatureContext
func newFeatureContext(sc *godog.ScenarioContext) *featureContext {
	fc := &featurecontext.FeatureContext{}

	// every xxxFeatureContext needs to have the pointer to a _single_ / common FeatureContext
	uc := &featureContext{
		FeatureContext: fc,

		LoginFeatureContext:       login.NewLoginFeatureContext(fc, sc),
		PublicShareFeatureContext: publicshare.NewPublicShareFeatureContext(fc, sc),
		ResourcesFeatureContext:   resources.NewResourcesFeatureContext(fc, sc),
		SpacesFeatureContext:      spaces.NewSpacesFeatureContext(fc, sc),
		SampleTestFeatureContext:  sampletest.NewSampleTestFeatureContext(fc, sc),
	}
	return uc
}
