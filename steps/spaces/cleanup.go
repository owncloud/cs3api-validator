package spaces

import (
	"context"

	"github.com/cucumber/godog"
)

// DeleteSpacesAfterScenario deletes all spaces which have been created after running the scenario
func (f *SpacesFeatureContext) DeleteSpacesAfterScenario(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	for range f.CreatedSpaces {

		//// deprovision space in the name of the user
		//reqctx, err := f.GetAuthContext(sp.Owner.Username)
		//if err != nil {
		//	return ctx, err
		//}

		// TODO: Deprovision storage spaces as soon as implemented
		//return nil, godog.ErrUndefined

		//resp, err := f.Client.DeleteStorageSpace(
		//	reqctx,
		//	&providerv1beta1.DeleteStorageSpaceRequest{
		//		Id: sp.Id,
		//	},
		//)
		//if err != nil {
		//	return ctx, err
		//}
		//if resp.Status.Code != rpc.Code_CODE_OK {
		//	return ctx, helpers.FormatError(resp.Status)
		//}
	}
	return ctx, nil
}
