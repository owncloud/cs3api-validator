package resources

import (
	"fmt"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/utils"
	"github.com/owncloud/cs3api-validator/helpers"
	"github.com/stretchr/testify/assert"
)

func (f *ResourcesFeatureContext) UserHasCreatedAResourceOfTypeInTheHomeDirectoryWithTheAlias(user, resourceName, resourceType, resourceAlias string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	homeSpace, err := f.GetHomeSpace(user)
	if err != nil {
		return err
	}

	statHome, err := f.Client.Stat(
		ctx,
		&providerv1beta1.StatRequest{
			Ref: &providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{StorageId: homeSpace.Root.StorageId, OpaqueId: homeSpace.Root.OpaqueId},
				Path: ".",
			},
		},
	)
	if err != nil {
		return err
	}
	if statHome.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(statHome.Status)
	}

	resourceRef := &providerv1beta1.Reference{
		ResourceId: statHome.Info.Id,
		Path:       utils.MakeRelativePath(resourceName),
	}

	switch resourceType {
	case "file":
		createResp, err := f.Client.TouchFile(
			ctx,
			&providerv1beta1.TouchFileRequest{
				Ref: resourceRef,
			},
		)
		if err != nil {
			return err
		}
		if createResp.Status.Code != rpc.Code_CODE_OK {
			return helpers.FormatError(createResp.Status)
		}

		f.Response = createResp

	case "container":
		createResp, err := f.Client.CreateContainer(
			ctx,
			&providerv1beta1.CreateContainerRequest{
				Ref: resourceRef,
			},
		)
		if err != nil {
			return err
		}
		if createResp.Status.Code != rpc.Code_CODE_OK {
			return helpers.FormatError(createResp.Status)
		}

		f.Response = createResp

	default:
		return fmt.Errorf("creating resource of type %s is not implemented", resourceType)
	}

	// store reference to delete on cleanup
	f.CreatedResourceReferences = append(f.CreatedResourceReferences, resourceRef)

	// store reference only if non empty alias
	if resourceAlias != "" {
		f.ResourceReferences[resourceAlias] = resourceRef
	}

	return nil
}

func (f *ResourcesFeatureContext) NoResourceShouldBeListedInTheResponse() error {
	list, ok := f.Response.(*providerv1beta1.ListContainerResponse)
	if !ok {
		return fmt.Errorf("expected to receive a ListContainerResponse but got something different")
	}

	return helpers.AssertExpectedAndActual(assert.Equal, 0, len(list.Infos))
}

func (f *ResourcesFeatureContext) ResourceOfTypeShouldBeListedInTheResponse(number int, resourceType string) error {
	list, ok := f.Response.(*providerv1beta1.ListContainerResponse)
	if !ok {
		return fmt.Errorf("expected to receive a ListContainerResponse but got something different")
	}

	var resType providerv1beta1.ResourceType
	switch resourceType {
	case "file":
		resType = providerv1beta1.ResourceType_RESOURCE_TYPE_FILE
	case "container":
		resType = providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER
	default:
		return fmt.Errorf("unknown resource type \"%s\"", resourceType)
	}

	matchingResources := 0

	for _, ri := range list.Infos {
		if ri.Type == resType {
			matchingResources++
		}
	}

	return helpers.AssertExpectedAndActual(assert.Equal, number, matchingResources)
}
