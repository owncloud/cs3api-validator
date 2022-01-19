package resources

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// ResourcesFeatureContext holds values which are used across test steps
type ResourcesFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewResourcesFeatureContext(fc *featurecontext.FeatureContext, sc *godog.ScenarioContext) *ResourcesFeatureContext {
	rfc := &ResourcesFeatureContext{FeatureContext: fc}
	rfc.Register(sc)
	return rfc
}

func (f *ResourcesFeatureContext) Register(sc *godog.ScenarioContext) {
	// steps
	sc.Step(`^no resource should be listed in the response$`, f.NoResourceShouldBeListedInTheResponse)
	sc.Step(`^(\d+) resource(?:s)? of type "([^"]*)" should be listed in the response$`, f.ResourceOfTypeShouldBeListedInTheResponse)
	sc.Step(`^user "([^"]*)" has created a folder "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.UserHasCreatedAFolderOfTypeInTheHomeDirectoryWithTheAlias)
	sc.Step(`^user "([^"]*)" has uploaded a file "([^"]*)" with content "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.userHasUploadedAFileWithContentInTheHomeDirectoryWithTheAlias)
	sc.Step(`^user "([^"]*)" remembers the etag of the space with name "([^"]*)"$`, f.userRemembersTheEtagOfTheSpaceWithName)

	// cleanup
	sc.After(f.DeleteResourcesAfterScenario)
	sc.After(f.EmptyTrashAfterScenario)
}
