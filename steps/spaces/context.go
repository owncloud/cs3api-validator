package spaces

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// SpacesFeatureContext holds values which are used across test steps
type SpacesFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewSpacesFeatureContext(fc *featurecontext.FeatureContext, sc *godog.ScenarioContext) *SpacesFeatureContext {
	spc := &SpacesFeatureContext{FeatureContext: fc}
	spc.Register(sc)
	return spc
}

func (f *SpacesFeatureContext) Register(ctx *godog.ScenarioContext) {
	// steps
	ctx.Step(`^user "([^"]*)" has created a personal space with the alias "([^"]*)"$`, f.UserHasCreatedAPersonalSpaceWithAlias)
	ctx.Step(`^user "([^"]*)" lists all available spaces$`, f.UserListsAllAvailableSpaces)
	ctx.Step(`^one personal space should be listed in the response$`, f.OnePersonalSpaceShouldBeListedInTheResponse)

	// cleanup
	ctx.After(f.DeleteSpacesAfterScenario)
}
