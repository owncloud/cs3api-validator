package main

import "github.com/cucumber/godog"

func aPersonalSpaceWasCreatedForUserAdmin() error {
	return godog.ErrPending
}

func iCallListStorageSpaces() error {
	return godog.ErrPending
}

func thereShouldBeOnePersonalSpaceInTheResponse() error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a personal space was created for user admin$`, aPersonalSpaceWasCreatedForUserAdmin)
	ctx.Step(`^I call listStorage Spaces$`, iCallListStorageSpaces)
	ctx.Step(`^there should be one personal space in the response$`, thereShouldBeOnePersonalSpaceInTheResponse)
}
