package login

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// LoginFeatureContext holds values which are used across test steps
type LoginFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewLoginFeatureContext(fc *featurecontext.FeatureContext, sc *godog.ScenarioContext) *LoginFeatureContext {
	lfc := &LoginFeatureContext{FeatureContext: fc}
	lfc.Register(sc)
	return lfc
}

func (f *LoginFeatureContext) Register(ctx *godog.ScenarioContext) {
	// steps
	ctx.Step(`^user "([^"]*)" has logged in with password "([^"]*)"$`, f.UserHasLoggedIn)
	ctx.Step(`^user "([^"]*)" has logged in with the token of the public-share "([^"]*)"$`, f.UserHasLoggedInWithTheTokenOfThePublicshare)

	// cleanup
}
