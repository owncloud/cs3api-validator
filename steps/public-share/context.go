package publicshare

import (
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// PublicShareFeatureContext holds values which are used across test steps
type PublicShareFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewPublicShareFeatureContext(fc *featurecontext.FeatureContext) *PublicShareFeatureContext {
	return &PublicShareFeatureContext{FeatureContext: fc}
}
