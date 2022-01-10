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

func (f *SpacesFeatureContext) Register(sc *godog.ScenarioContext){
	// steps
	sc.Step(`^user "([^"]*)" has created a personal space$`, f.UserHasCreatedAPersonalSpace)
	sc.Step(`^user "([^"]*)" lists all available spaces$`, f.UserListsAllAvailableSpaces)
	sc.Step(`^one personal space should be listed in the response$`, f.OnePersonalSpaceShouldBeListedInTheResponse)

	// cleanup
	sc.After(f.DeleteSpacesAfterScenario)
}
