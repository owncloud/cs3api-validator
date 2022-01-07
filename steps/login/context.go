package login

import (
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// LoginFeatureContext holds values which are used across test steps
type LoginFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewLoginFeatureContext(fc *featurecontext.FeatureContext) *LoginFeatureContext {
	return &LoginFeatureContext{FeatureContext: fc}
}
