package resources

import "github.com/owncloud/cs3api-validator/featurecontext"

// ResourcesFeatureContext holds values which are used across test steps
type ResourcesFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewResourcesFeatureContext(fc *featurecontext.FeatureContext) *ResourcesFeatureContext {
	return &ResourcesFeatureContext{FeatureContext: fc}
}
