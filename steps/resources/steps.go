package resources

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"
	"github.com/stretchr/testify/assert"
)

const (
	// HeaderTokenTransport holds the header key for the reva transfer token
	// "github.com/cs3org/reva/v2/internal/http/services/datagateway" is internal so we redeclare it here
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
				ResourceId: homeSpace.Root,
				Path:       ".",
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
				ResourceId: homeSpace.Root,
				Path:       ".",
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

	length := int64(len(content))
	uReq := &providerv1beta1.InitiateFileUploadRequest{
		Opaque: utils.AppendPlainToOpaque(nil, net.HeaderUploadLength, strconv.FormatInt(length, 10)),
		Ref:    resourceRef,
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

	if length > 0 {
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

		ok := chunking.IsChunked(resourceRef.Path)
		if ok {
			chunk, err := chunking.GetChunkBLOBInfo(resourceRef.Path)
			if err != nil {
				return err
			}
			sReq = &providerv1beta1.StatRequest{Ref: &providerv1beta1.Reference{Path: chunk.Path}}
		}
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

func (f *ResourcesFeatureContext) theFollowingResourcesShouldBeListedInTheResponse(not string, table *godog.Table) error {
	list, ok := f.Response.(*providerv1beta1.ListContainerResponse)
	if !ok {
		return fmt.Errorf("expected to receive a ListContainerResponse but got something different")
	}
	rows := table.Rows
	if len(rows) == 0 {
		return fmt.Errorf("empty gherkin table")
	}
	if rows[0].Cells[0].Value != "type" && rows[0].Cells[1].Value != "path" {
		return fmt.Errorf("the first line of the tables needs to be in the form | <type> | <path> |")
	}
	var rowsValues []*messages.PickleTableRow
	// collect row values and leave the table header row out
	rowsValues = append(rowsValues, rows[+1:]...)
	matchingResources := make(map[string]bool)
	// loop over container resources
	for _, ri := range list.Infos {
		// check if one of expected values has a matching entry in the container resources
		for _, row := range rowsValues {
			var resType providerv1beta1.ResourceType
			switch row.Cells[0].Value {
			case "file":
				resType = providerv1beta1.ResourceType_RESOURCE_TYPE_FILE
			case "container":
				resType = providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER
			default:
				return fmt.Errorf("unknown resource type \"%s\"", row.Cells[0].Value)
			}
			expectedPath := row.Cells[1].Value
			if resType == ri.Type && expectedPath == ri.Path {
				matchingResources[row.Cells[1].Value] = true
				if not == "not" {
					msg := fmt.Sprintf("Resource with path %s should not be listed on the response", expectedPath)
					return helpers.AssertActual(assert.Zero, matchingResources, msg)
				}
			}
		}
	}
	if not == "" {
		for _, candidate := range rowsValues {
			path := candidate.Cells[1].Value
			err := helpers.AssertExpectedAndActual(assert.Equal, matchingResources[path], true)
			if err != nil {
				return fmt.Errorf("the resource with path %s could not be found in the response", path)
			}
		}
	}
	return nil
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

func (f *ResourcesFeatureContext) forUserTheChecksumsOfTheResourceWithTheAliasShouldHaveChanged(user string, alias string, not string) error {
	reqctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	resource, ok := f.ResourceReferences[alias]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", alias)
	}
	if resource.Info.Type != providerv1beta1.ResourceType_RESOURCE_TYPE_FILE {
		return fmt.Errorf("checksums are only available on files")
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

	return helpers.AssertExpectedAndActual(assertion, resource.Info.Checksum.GetSum(), resp.Info.Checksum.GetSum())
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

func (f *ResourcesFeatureContext) userMovesTheResourceWithAliasInsideASpaceToTarget(user string, alias string, target string) error {
	reqctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	// get resource from alias
	resource, ok := f.ResourceReferences[alias]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", alias)
	}
	// create target reference with new path
	targetRef := providerv1beta1.Reference{
		ResourceId: resource.Ref.ResourceId,
		Path:       utils.MakeRelativePath(target),
	}
	// check if target already exists
	tStat, err := f.Client.Stat(
		reqctx,
		&providerv1beta1.StatRequest{
			Ref: &targetRef,
		},
	)
	if err != nil {
		return err
	}

	// abort if target already exists
	if tStat.Status.Code != rpc.Code_CODE_NOT_FOUND {
		if tStat.Status.Code == rpc.Code_CODE_OK {
			return fmt.Errorf("Resource already exists")
		}
		return helpers.FormatError(tStat.Status)
	}
	// do the move
	mRes, err := f.Client.Move(
		reqctx, &providerv1beta1.MoveRequest{
			Source:      resource.Ref,
			Destination: &targetRef,
		},
	)
	if err != nil {
		return err
	}
	if mRes.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(mRes.Status)
	}

	// stat the new location
	sResp, err := f.Client.Stat(reqctx, &providerv1beta1.StatRequest{
		Ref: &targetRef,
	})
	if err != nil {
		return err
	}
	if sResp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(sResp.Status)
	}
	// store reference to delete on cleanup
	// we only add the new references, the outdated ones can remain
	f.CreatedResourceReferences = append(f.CreatedResourceReferences, &targetRef)

	// update the remembered references
	f.ResourceReferences[alias] = featurecontext.ResourceAlias{Ref: &targetRef, Info: sResp.Info}
	// update the remembered child aliases
	return f.updateChildAliases(reqctx, f.ResourceReferences[alias])
}

// updateChildAliases makes sure that child nodes are up-to-date in the aliases index
func (f *ResourcesFeatureContext) updateChildAliases(ctx context.Context, parent featurecontext.ResourceAlias) error {
	if parent.Info.Type != providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER {
		return nil
	}
	// use a stack to explore sub-containers breadth-first
	stack := []string{parent.Ref.Path}
	for len(stack) > 0 {
		// retrieve path on top of stack
		path := stack[len(stack)-1]

		nRef := &providerv1beta1.Reference{
			ResourceId: parent.Ref.ResourceId,
			Path:       path,
		}
		req := &providerv1beta1.ListContainerRequest{
			Ref: nRef,
		}
		res, err := f.Client.ListContainer(ctx, req)
		if err != nil {
			return fmt.Errorf("transport error sending list container grpc request")
		}
		if res.Status.Code != rpc.Code_CODE_OK {
			return fmt.Errorf("error sending list container grpc request")
		}

		stack = stack[:len(stack)-1]

		// check sub-containers in reverse order and add them to the stack
		// the reversed order here will produce a more logical sorting of results
		for i := len(res.Infos) - 1; i >= 0; i-- {
			res.Infos[i].Path = utils.MakeRelativePath(filepath.Join(nRef.Path, res.Infos[i].Path))
			currentRef := providerv1beta1.Reference{
				ResourceId: parent.Ref.ResourceId,
				Path:       res.Infos[i].Path,
			}
			resourceInfo := featurecontext.ResourceAlias{
				Ref:  &currentRef,
				Info: res.Infos[i],
			}
			if res.Infos[i].Type == providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER {
				stack = append(stack, res.Infos[i].Path)
			}
			// add to existing references
			f.CreatedResourceReferences = append(f.CreatedResourceReferences, &currentRef)
			// update the existing alias
			for alias, ri := range f.ResourceReferences {
				if ri.Info.Id.OpaqueId == res.Infos[i].Id.OpaqueId {
					f.ResourceReferences[alias] = resourceInfo
				}
			}
		}
	}
	return nil
}

func (f *ResourcesFeatureContext) userListsAllResourcesInsideTheResourceWithAlias(user string, alias string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}
	var res featurecontext.ResourceAlias
	res, ok := f.ResourceReferences[alias]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", alias)
	}
	if res.Info.Type != providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER {
		return fmt.Errorf("we cannot call list inside non-container resources")
	}

	resp, err := f.Client.ListContainer(
		ctx,
		&providerv1beta1.ListContainerRequest{
			Ref: res.Ref,
		},
	)
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}

	f.Response = resp

	return nil
}
