package spaces

import (
	"errors"
	"fmt"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/cs3api-validator/helpers"
	"github.com/stretchr/testify/assert"
)

func (f *SpacesFeatureContext) UserHasCreatedAPersonalSpace(user string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

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
	filters = append(filters, filterHome)
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
	return nil
}

func (f *SpacesFeatureContext) UserListsAllAvailableSpaces(user string) error {
	ctx, err := f.GetAuthContext(user)
	if err != nil {
		return err
	}

	resp, err := f.Client.ListStorageSpaces(
		ctx,
		&providerv1beta1.ListStorageSpacesRequest{
			Filters: nil,
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
