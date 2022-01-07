package resources

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// ResourcesFeatureContext holds values which are used across test steps
type ResourcesFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewResourcesFeatureContext(fc *featurecontext.FeatureContext) *ResourcesFeatureContext {
	return &ResourcesFeatureContext{FeatureContext: fc}
}

func (f *ResourcesFeatureContext) RegisterSteps(sc *godog.ScenarioContext) {
	// steps
	sc.Step(`^no resource should be listed in the response$`, f.NoResourceShouldBeListedInTheResponse)
	sc.Step(`^(\d+) resource(?:s)? of type "([^"]*)" should be listed in the response$`, f.ResourceOfTypeShouldBeListedInTheResponse)
	sc.Step(`^user "([^"]*)" has created a resource "([^"]*)" of type "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.UserHasCreatedAResourceOfTypeInTheHomeDirectoryWithTheAlias)

	// cleanup
	sc.After(f.DeleteResourcesAfterScenario)
	sc.After(f.EmptyTrashAfterScenario)
}
