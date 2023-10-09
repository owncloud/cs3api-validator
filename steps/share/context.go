package share

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// ShareTestFeatureContext holds values which are used across test steps
type ShareTestFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewShareTestFeatureContext(fc *featurecontext.FeatureContext, sc *godog.ScenarioContext) *ShareTestFeatureContext {
	nsc := &ShareTestFeatureContext{FeatureContext: fc}
	nsc.Register(sc)
	return nsc
}

func (f *ShareTestFeatureContext) Register(ctx *godog.ScenarioContext) {
	// steps
	ctx.Step(`^user "([^"]*)" shares a file "([^"]*)" with the following users concurrently$`, f.UserSharesAFileWithTheFollowingUsers)
	ctx.Step(`^the concurrent user sharing should have been successfull$`, f.TheConcurrentUserSharingShouldHaveBeenSuccessfull)
}
