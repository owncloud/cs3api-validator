package scenario

import (
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/steps/login"
	publicshare "github.com/owncloud/cs3api-validator/steps/public-share"
	"github.com/owncloud/cs3api-validator/steps/resources"
	"github.com/owncloud/cs3api-validator/steps/spaces"
)

// FeatureContext holds values which are used across test steps
type featureContext struct {
	*featurecontext.FeatureContext

	*login.LoginFeatureContext
	*publicshare.PublicShareFeatureContext
	*resources.ResourcesFeatureContext
	*spaces.SpacesFeatureContext
}

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
