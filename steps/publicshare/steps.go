package publicshare

import (
	"fmt"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	linkv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/cs3api-validator/helpers"
)

func (f *PublicShareFeatureContext) UserHasCreatedAPublicshareWithEditorPermissionsOfTheResourceWithTheAlias(user, publicshare, resourceAlias string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	resourceRef, ok := f.ResourceReferences[resourceAlias]
	if !ok {
		return fmt.Errorf("resource alias %s is not known", resourceAlias)
	}

	statResp, err := f.Client.Stat(
		ctx,
		&providerv1beta1.StatRequest{
			Ref: resourceRef.Ref,
		},
	)
	if err != nil {
		return err
	}
	if statResp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(statResp.Status)
	}

	resp, err := f.Client.CreatePublicShare(
		ctx,
		&linkv1beta1.CreatePublicShareRequest{
			// TODO: how to set the display name of a public share?
			ResourceInfo: statResp.Info,
			Grant: &linkv1beta1.Grant{
				Permissions: &linkv1beta1.PublicSharePermissions{
					Permissions: &providerv1beta1.ResourcePermissions{
						// see https://github.com/cs3org/reva/blob/master/cmd/reva/share-create.go#L160-L204
						GetPath:              true,
						InitiateFileDownload: true,
						ListFileVersions:     true,
						ListContainer:        true,
						Stat:                 true,
						CreateContainer:      true,
						Delete:               true,
						InitiateFileUpload:   true,
						RestoreFileVersion:   true,
						Move:                 true,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}

	f.PublicSharesToken[publicshare] = resp.Share.Token

	f.Response = resp

	return err

}

func (f *PublicShareFeatureContext) UserListsAllResourcesInThePublicshare(user, publicShare string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	token, err := f.GetPublicShareToken(publicShare)
	if err != nil {
		return err
	}

	publicShareResp, err := f.Client.GetPublicShareByToken(
		ctx,
		&linkv1beta1.GetPublicShareByTokenRequest{
			Token: token,
		},
	)

	if err != nil {
		return err
	}
	if publicShareResp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(publicShareResp.Status)
	}

	resp, err := f.Client.ListContainer(
		ctx,
		&providerv1beta1.ListContainerRequest{
			Ref: &providerv1beta1.Reference{
				ResourceId: publicShareResp.Share.ResourceId,
			},
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

func (f *PublicShareFeatureContext) UserHasUploadedAnEmptyFileToThePublicshare(user, filename, publicShare string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	token, err := f.GetPublicShareToken(publicShare)
	if err != nil {
		return err
	}

	publicShareResp, err := f.Client.GetPublicShareByToken(
		ctx,
		&linkv1beta1.GetPublicShareByTokenRequest{
			Token: token,
		},
	)

	if err != nil {
		return err
	}
	if publicShareResp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(publicShareResp.Status)
	}

	// TODO: switch to TouchFile when implemented in REVA
	//resp, err := f.Client.TouchFile(
	//	ctx,
	//	&providerv1beta1.TouchFileRequest{
	//		Ref: &providerv1beta1.Reference{
	//			ResourceId: share.ResourceId,
	//			Path:       utils.MakeRelativePath(filename),
	//		},
	//	},
	//)

	// TODO: how can one access the share through the public storage provider?
	//publicShareResp.Share.ResourceId.StorageId = "e1a73ede-549b-4226-abdf-40e69ca8230d"

	// TODO: how to do this without magic?
	path := "/public/" + token + "/" + filename

	resp, err := f.Client.InitiateFileUpload(
		ctx,
		&providerv1beta1.InitiateFileUploadRequest{
			Ref: &providerv1beta1.Reference{
				//ResourceId: publicShareResp.Share.ResourceId,
				Path: path,
			},
		},
	)
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}

	f.Response = nil
	return nil
}
