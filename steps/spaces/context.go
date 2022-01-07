package spaces

import "github.com/owncloud/cs3api-validator/featurecontext"

// SpacesFeatureContext holds values which are used across test steps
type SpacesFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewSpacesFeatureContext(fc *featurecontext.FeatureContext) *SpacesFeatureContext {
	return &SpacesFeatureContext{FeatureContext: fc}
}
