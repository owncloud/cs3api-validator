package spaces

import (
	"encoding/json"
	"errors"
	"fmt"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"
	"github.com/stretchr/testify/assert"
)

func (f *SpacesFeatureContext) UserHasCreatedAPersonalSpaceWithAlias(user string, alias string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	// Call create home
	resp, err := f.Client.CreateHome(
		ctx,
		&providerv1beta1.CreateHomeRequest{},
	)
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK && resp.Status.Code != rpc.Code_CODE_ALREADY_EXISTS {
		return helpers.FormatError(resp.Status)
	}
	var filters []*providerv1beta1.ListStorageSpacesRequest_Filter
	filterHome := &providerv1beta1.ListStorageSpacesRequest_Filter{
		Type: providerv1beta1.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
		Term: &providerv1beta1.ListStorageSpacesRequest_Filter_SpaceType{
			SpaceType: "personal",
		},
	}
	filterUser := &providerv1beta1.ListStorageSpacesRequest_Filter{
		Type: providerv1beta1.ListStorageSpacesRequest_Filter_TYPE_OWNER,
		Term: &providerv1beta1.ListStorageSpacesRequest_Filter_Owner{
			Owner: f.Users[user].User.Id,
		},
	}
	filters = append(filters, filterHome, filterUser)
	homeResp, err := f.Client.ListStorageSpaces(
		ctx,
		&providerv1beta1.ListStorageSpacesRequest{
			Filters: filters,
		},
	)
	if err != nil {
		return err
	}
	f.Response = homeResp
	personalRef := &providerv1beta1.Reference{
		ResourceId: homeResp.StorageSpaces[0].Root,
		Path:       ".",
	}
	statPersonal, err := f.Client.Stat(
		ctx,
		&providerv1beta1.StatRequest{Ref: personalRef},
	)
	if err != nil {
		return err
	}
	if statPersonal.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(statPersonal.Status)
	}
	// store reference to delete on cleanup
	f.CreatedSpaces = append(f.CreatedSpaces, homeResp.StorageSpaces[0])

	// store reference only if non-empty alias
	if alias != "" {
		f.ResourceReferences[alias] = featurecontext.ResourceAlias{
			Ref: &providerv1beta1.Reference{
				ResourceId: homeResp.StorageSpaces[0].Root,
				Path:       ".",
			},
			Info: statPersonal.Info,
		}
	}
	return nil
}

func (f *SpacesFeatureContext) UserListsAllAvailableSpaces(user string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}
	// we need to send empty permissions as a workaround
	permissions := make(map[string]struct{}, 1)
	value, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	resp, err := f.Client.ListStorageSpaces(
		ctx,
		&providerv1beta1.ListStorageSpacesRequest{
			Opaque: &types.Opaque{Map: map[string]*types.OpaqueEntry{
				"permissions": {
					Decoder: "json",
					Value:   value,
				},
			}},
		},
	)
	if err != nil {
		return err
	}
	if resp.Status.Code != rpc.Code_CODE_OK {
		return helpers.FormatError(resp.Status)
	}

	f.Response = resp

	return err

}

func (f *SpacesFeatureContext) OnePersonalSpaceShouldBeListedInTheResponse() error {
	if f.Response == nil {
		return errors.New("no response available")
	}

	resp, ok := f.Response.(*providerv1beta1.ListStorageSpacesResponse)
	if !ok {
		return fmt.Errorf("no valid response from former requests available: %v", resp)
	}

	var personalSpaces []*providerv1beta1.StorageSpace
	for _, s := range resp.StorageSpaces {
		if s.SpaceType == "personal" {
			personalSpaces = append(personalSpaces, s)
		}
	}
	return helpers.AssertExpectedAndActual(assert.Equal, 1, len(personalSpaces))
}
