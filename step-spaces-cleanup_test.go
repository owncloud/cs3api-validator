package main

import (
	"context"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cucumber/godog"
)

// deleteSpacesAfterScenario deletes all spaces which have been created after running the scenario
func (f *FeatureContext) deleteSpacesAfterScenario(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	for _, sp := range f.CreatedSpaces {

		// deprovision space in the name of the user
		reqctx, err := f.getAuthContext(sp.Owner.Username)
		if err != nil {
			return ctx, err
		}

		// TODO: Deprovision storage spaces as soon as implemented
		return nil, godog.ErrUndefined

		resp, err := f.Client.DeleteStorageSpace(
			reqctx,
			&providerv1beta1.DeleteStorageSpaceRequest{
				Id: sp.Id,
			},
		)
		if err != nil {
			return ctx, err
		}
		if resp.Status.Code != rpc.Code_CODE_OK {
			return ctx, formatError(resp.Status)
		}

	}
	return ctx, nil
}
