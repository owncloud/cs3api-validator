package sampletest

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// SampleTestFeatureContext holds values which are used across test steps
type SampleTestFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewSampleTestFeatureContext(fc *featurecontext.FeatureContext, sc *godog.ScenarioContext) *SampleTestFeatureContext {
	nsc := &SampleTestFeatureContext{FeatureContext: fc}
	nsc.Register(sc)
	return nsc
}

func (f *SampleTestFeatureContext) Register(ctx *godog.ScenarioContext) {
	// steps
	ctx.Step(`^user "([^"]*)" shares a file "([^"]*)" with the following users concurrently$`, f.userSharesAFileWithTheFollowingUsers)
}
