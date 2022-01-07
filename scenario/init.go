package scenario

import (
	"github.com/cucumber/godog"
)

var Endpoint string
var HttpInsecure bool

func InitializeScenario(endpoint string, httpInsecure bool) func(*godog.ScenarioContext) {

	Endpoint = endpoint
	HttpInsecure = httpInsecure

	return initializeScenario
}

func initializeScenario(sc *godog.ScenarioContext) {

	f := newFeatureContext()

	f.Init(
		Endpoint,
		HttpInsecure,
	)

	// Deprovision all storage spaces after the scenario
	sc.After(f.DeletePublicShares)
	sc.After(f.DeleteSpacesAfterScenario)
	sc.After(f.DeleteResourcesAfterScenario)
	sc.After(f.EmptyTrashAfterScenario)

	// Step implementations

	// login
	sc.Step(`^user "([^"]*)" has logged in with password "([^"]*)"$`, f.UserHasLoggedIn)
	sc.Step(`^user "([^"]*)" has logged in with the token of the public-share "([^"]*)"$`, f.UserHasLoggedInWithTheTokenOfThePublicshare)

	// spaces
	sc.Step(`^user "([^"]*)" has created a personal space$`, f.UserHasCreatedAPersonalSpace)
	sc.Step(`^user "([^"]*)" lists all available spaces$`, f.UserListsAllAvailableSpaces)
	sc.Step(`^one personal space should be listed in the response$`, f.OnePersonalSpaceShouldBeListedInTheResponse)

	// create resources
	sc.Step(`^user "([^"]*)" has created a resource "([^"]*)" of type "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.UserHasCreatedAResourceOfTypeInTheHomeDirectoryWithTheAlias)

	// list resources
	sc.Step(`^user "([^"]*)" lists all resources in the public-share "([^"]*)"$`, f.UserListsAllResourcesInThePublicshare)
	sc.Step(`^no resource should be listed in the response$`, f.NoResourceShouldBeListedInTheResponse)
	sc.Step(`^(\d+) resource(?:s)? of type "([^"]*)" should be listed in the response$`, f.ResourceOfTypeShouldBeListedInTheResponse)

	// public shares
	sc.Step(`^user "([^"]*)" has created a public-share "([^"]*)" with editor permissions of the resource with the alias "([^"]*)"$`, f.UserHasCreatedAPublicshareWithEditorPermissionsOfTheResourceWithTheAlias)
	sc.Step(`^user "([^"]*)" has uploaded an empty file "([^"]*)" to the public-share "([^"]*)"$`, f.UserHasUploadedAnEmptyFileToThePublicshare)

}
