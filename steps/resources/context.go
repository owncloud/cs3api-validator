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

func (f *ResourcesFeatureContext) Register(ctx *godog.ScenarioContext) {
	// steps
	ctx.Step(`^no resource should be listed in the response$`, f.NoResourceShouldBeListedInTheResponse)
	ctx.Step(`^(\d+) resource(?:s)? of type "([^"]*)" should be listed in the response$`, f.ResourceOfTypeShouldBeListedInTheResponse)
	ctx.Step(`^user "([^"]*)" has created a folder "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.UserHasCreatedAFolderOfTypeInTheHomeDirectoryWithTheAlias)
	ctx.Step(`^user "([^"]*)" has uploaded a file "([^"]*)" with content "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.userHasUploadedAFileWithContentInTheHomeDirectoryWithTheAlias)
	ctx.Step(`^user "([^"]*)" remembers the fileinfo of the resource with the alias "([^"]*)"$`, f.userRemembersTheFileInfoOfTheResourceWithTheAlias)
	ctx.Step(`^for user "([^"]*)" the etag of the resource with the alias "([^"]*)" should (not|)\s?have changed$`, f.forUserTheEtagOfTheResourceWithTheAliasShouldHaveChanged)
	ctx.Step(`^for user "([^"]*)" the treesize of the resource with the alias "([^"]*)" should be (\d+)$`, f.forUserTheTreesizeOfTheResourceWithTheAliasShouldBe)

	// cleanup
	ctx.After(f.DeleteResourcesAfterScenario)
	ctx.After(f.EmptyTrashAfterScenario)
}
