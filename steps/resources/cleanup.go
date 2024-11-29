package resources

import (
	"context"
	"errors"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cucumber/godog"
	"github.com/owncloud/cs3api-validator/helpers"
)

// deleteResourcesAfterScenario deletes resources which have been created after running the scenario
func (f *ResourcesFeatureContext) DeleteResourcesAfterScenario(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {

	resourcesToDelete := f.CreatedResourceReferences

	// we don't now which user has access to which reference,
	// therefore we just try to delete each reference with all users
	for u := range f.Users {
		reqctx, err := f.GetAuthContext(u)
		if err != nil {
			continue
		}

		notYetDeleted := []*providerv1beta1.Reference{}

		for _, ref := range resourcesToDelete {
			// we need to retry the deletion because the resource might be in postprocessing
			e := helpers.RetryResource(f.deleteResource, ref, reqctx, 3, f.Config.AsyncPropagationDelay)
			if e != nil {
				notYetDeleted = append(notYetDeleted, ref)
			}
		}

		resourcesToDelete = notYetDeleted

	}
	return ctx, nil
}

func (f *ResourcesFeatureContext) deleteResource(ref *providerv1beta1.Reference, reqctx context.Context) error {
	resp, err := f.Client.Delete(
		reqctx,
		&providerv1beta1.DeleteRequest{
			Ref: ref,
		},
	)
	if err != nil {
		return err
	}

	// non-existing resources are not errors
	if resp.Status.Code != rpc.Code_CODE_OK && resp.Status.Code != rpc.Code_CODE_NOT_FOUND {
		return errors.New(resp.GetStatus().GetMessage())
	}
	return nil
}

// emptyTrashAfterScenario empties the trash for all users after running the scenario
func (f *ResourcesFeatureContext) EmptyTrashAfterScenario(ctx context.Context, sc *godog.Scenario, errsctx error) (context.Context, error) {
	for u := range f.Users {
		if u == "anonymous" {
			continue
		}
		reqctx, err := f.GetAuthContext(u)
		if err != nil {
			continue
		}

		homeSpace, err := f.GetHomeSpace(u)

		if err != nil {
			continue
		}

		//TODO: oCIS FS -> why do we still have blobs on disk?
		_, _ = f.Client.PurgeRecycle(
			reqctx,
			&providerv1beta1.PurgeRecycleRequest{
				Ref: &providerv1beta1.Reference{
					ResourceId: &providerv1beta1.ResourceId{OpaqueId: homeSpace.Root.OpaqueId, StorageId: homeSpace.Root.OpaqueId},
					Path:       ".",
				},
			},
		)

	}
	return ctx, nil
}
