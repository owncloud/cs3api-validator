package main

import (
	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {

	f := &FeatureContext{}
	f.Init()

	// Deprovision all storage spaces after the scenario
	ctx.After(f.deleteSpacesAfterScenario)
	ctx.After(f.deleteResourcesAfterScenario)
	ctx.After(f.emptyTrashAfterScenario)

	// Step implementations

	// login
	ctx.Step(`^user "([^"]*)" has logged in with password "([^"]*)"$`, f.userHasLoggedIn)
	ctx.Step(`^user "([^"]*)" has logged in with the token of the public-share "([^"]*)"$`, f.userHasLoggedInWithTheTokenOfThePublicshare)

	// spaces
	ctx.Step(`^user "([^"]*)" has created a personal space$`, f.userHasCreatedAPersonalSpace)
	ctx.Step(`^user "([^"]*)" lists all available spaces$`, f.userListsAllAvailableSpaces)
	ctx.Step(`^one personal space should be listed in the response$`, f.onePersonalSpaceShouldBeListedInTheResponse)

	// create resources
	ctx.Step(`^user "([^"]*)" has created a resource "([^"]*)" of type "([^"]*)" in the home directory with the alias "([^"]*)"$`, f.userHasCreatedAResourceOfTypeInTheHomeDirectoryWithTheAlias)

	// list resources
	ctx.Step(`^user "([^"]*)" lists all resources in the public-share "([^"]*)"$`, f.userListsAllResourcesInThePublicshare)
	ctx.Step(`^no resource should be listed in the response$`, f.noResourceShouldBeListedInTheResponse)
	ctx.Step(`^(\d+) resource(?:s)? of type "([^"]*)" should be listed in the response$`, f.resourceOfTypeShouldBeListedInTheResponse)

	// public shares
	ctx.Step(`^user "([^"]*)" has created a public-share "([^"]*)" with editor permissions of the resource with the alias "([^"]*)"$`, f.userHasCreatedAPublicshareWithEditorPermissionsOfTheResourceWithTheAlias)
	ctx.Step(`^user "([^"]*)" has uploaded an empty file "([^"]*)" to the public-share "([^"]*)"$`, f.userHasUploadedAnEmptyFileToThePublicshare)

}
