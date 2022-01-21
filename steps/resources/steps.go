package resources

import (
	"fmt"
	"net/http"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/errtypes"
	"github.com/cs3org/reva/pkg/rhttp"
	"github.com/cs3org/reva/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/pkg/utils"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"
	"github.com/stretchr/testify/assert"
)

const (
	// HeaderTokenTransport holds the header key for the reva transfer token
	// "github.com/cs3org/reva/internal/http/services/datagateway" is internal so we redeclare it here
	HeaderTokenTransport = "X-Reva-Transfer"
)

func (f *ResourcesFeatureContext) UserHasCreatedAFolderOfTypeInTheHomeDirectoryWithTheAlias(user, resourceName, resourceAlias string) error {
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
				ResourceId: &providerv1beta1.ResourceId{
					OpaqueId:  homeSpace.Root.OpaqueId,
					StorageId: homeSpace.Root.StorageId,
				},
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

	sRes, err := f.Client.Stat(
		ctx,
		&providerv1beta1.StatRequest{
			Ref: resourceRef,
		},
	)

	if err != nil {
		return err
	}

	if sRes.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("error statting new folder")
	}

	// store reference to delete on cleanup
	f.CreatedResourceReferences = append(f.CreatedResourceReferences, resourceRef)

	// store reference only if non empty alias
	if resourceAlias != "" {
		f.ResourceReferences[resourceAlias] = featurecontext.ResourceAlias{
			Ref:  resourceRef,
			Info: sRes.Info,
		}
	}

	return nil
}

func (f *ResourcesFeatureContext) userHasUploadedAFileWithContentInTheHomeDirectoryWithTheAlias(user, resourceName, content, resourceAlias string) error {
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
				ResourceId: &providerv1beta1.ResourceId{
					OpaqueId:  homeSpace.Root.OpaqueId,
					StorageId: homeSpace.Root.StorageId,
				},
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
	sReq := &providerv1beta1.StatRequest{Ref: resourceRef}
	sRes, err := f.Client.Stat(ctx, sReq)
	if err != nil {
		return err
	}
	if sRes.Status.Code != rpc.Code_CODE_OK && sRes.Status.Code != rpc.Code_CODE_NOT_FOUND {
		return fmt.Errorf("could not stat upload target")
	}

	info := sRes.Info
	if info != nil {
		if info.Type != providerv1beta1.ResourceType_RESOURCE_TYPE_FILE {
			return fmt.Errorf("target resource is not a file")
		}
	}

	uReq := &providerv1beta1.InitiateFileUploadRequest{
		Ref: resourceRef,
	}

	// where to upload the file?
	uRes, err := f.Client.InitiateFileUpload(ctx, uReq)
	if err != nil {
		return err
	}
	if uRes.Status.Code != rpc.Code_CODE_OK {
		switch uRes.Status.Code {
		case rpc.Code_CODE_PERMISSION_DENIED:
			return fmt.Errorf("permission denied to initiate upload %v", uReq)
		case rpc.Code_CODE_NOT_FOUND:
			return fmt.Errorf("target not found to initiate upload %v", uReq)
		default:
			return fmt.Errorf("error occurred during upload %v", uReq)
		}
	}

	var ep, token string
	for _, p := range uRes.Protocols {
		if p.Protocol == "simple" {
			ep, token = p.UploadEndpoint, p.Token
		}
	}
	body := strings.NewReader(content)
	httpReq, err := rhttp.NewRequest(ctx, http.MethodPut, ep, body)
	if err != nil {
		return err
	}
	httpReq.Header.Set(HeaderTokenTransport, token)

	httpRes, err := f.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	if httpRes.StatusCode != http.StatusOK {
		if httpRes.StatusCode == http.StatusPartialContent {
			return fmt.Errorf("partial content")
		}
		if httpRes.StatusCode == errtypes.StatusChecksumMismatch {
			return fmt.Errorf("checksum mismatch")
		}
		return fmt.Errorf("PUT request to datagateway failed")
	}

	ok, err := chunking.IsChunked(resourceRef.Path)
	if err != nil {
		return err
	}
	if ok {
		chunk, err := chunking.GetChunkBLOBInfo(resourceRef.Path)
		if err != nil {
			return err
		}
		sReq = &providerv1beta1.StatRequest{Ref: &providerv1beta1.Reference{Path: chunk.Path}}
	}

	// stat again to check the new file's metadata
	sRes, err = f.Client.Stat(ctx, sReq)
	if err != nil {
		return err
	}

	if sRes.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("error statting new file")
	}

	newInfo := sRes.Info

	// store reference to delete on cleanup
	f.CreatedResourceReferences = append(f.CreatedResourceReferences, resourceRef)

	// store reference only if non-empty alias
	if resourceAlias != "" {
		f.ResourceReferences[resourceAlias] = featurecontext.ResourceAlias{Ref: resourceRef, Info: newInfo}
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

func (f *ResourcesFeatureContext) userRemembersTheFileInfoOfTheResourceWithTheAlias(user string, alias string) error {
	reqctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	resource, ok := f.ResourceReferences[alias]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", alias)
	}

	resp, err := f.Client.Stat(reqctx, &providerv1beta1.StatRequest{
		Ref: resource.Ref,
	})
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}
	f.ResourceReferences[alias] = featurecontext.ResourceAlias{Ref: resource.Ref, Info: resp.Info}
	return nil
}

func (f *ResourcesFeatureContext) forUserTheEtagOfTheResourceWithTheAliasShouldHaveChanged(user string, alias string, not string) error {
	reqctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	resource, ok := f.ResourceReferences[alias]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", alias)
	}
	resp, err := f.Client.Stat(reqctx, &providerv1beta1.StatRequest{
		Ref: resource.Ref,
	})
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}

	var assertion helpers.ExpectedAndActualAssertion
	switch not {
	// not have changed is equal
	case "not":
		assertion = assert.Equal
	default:
		assertion = assert.NotEqual
	}

	return helpers.AssertExpectedAndActual(assertion, resource.Info.Etag, resp.Info.Etag)
}

func (f *ResourcesFeatureContext) forUserTheTreesizeOfTheResourceWithTheAliasShouldBe(user string, alias string, size int) error {
	reqctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	resource, ok := f.ResourceReferences[alias]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", alias)
	}
	resp, err := f.Client.Stat(reqctx, &providerv1beta1.StatRequest{
		Ref: resource.Ref,
	})
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}

	return helpers.AssertExpectedAndActual(assert.Equal, resp.Info.Size, uint64(size))
}
