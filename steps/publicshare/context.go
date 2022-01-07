package publicshare

import (
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/featurecontext"
)

// PublicShareFeatureContext holds values which are used across test steps
type PublicShareFeatureContext struct {
	*featurecontext.FeatureContext
}

func NewPublicShareFeatureContext(fc *featurecontext.FeatureContext) *PublicShareFeatureContext {
	return &PublicShareFeatureContext{FeatureContext: fc}
}

func (f *PublicShareFeatureContext) RegisterSteps(sc *godog.ScenarioContext) {
	// steps
	sc.Step(`^user "([^"]*)" has created a public-share "([^"]*)" with editor permissions of the resource with the alias "([^"]*)"$`, f.UserHasCreatedAPublicshareWithEditorPermissionsOfTheResourceWithTheAlias)
	sc.Step(`^user "([^"]*)" has uploaded an empty file "([^"]*)" to the public-share "([^"]*)"$`, f.UserHasUploadedAnEmptyFileToThePublicshare)
	sc.Step(`^user "([^"]*)" lists all resources in the public-share "([^"]*)"$`, f.UserListsAllResourcesInThePublicshare)

	// cleanup
	sc.After(f.DeletePublicShares)
}
