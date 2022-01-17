package featurecontext

import (
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

// GetHomeSpace finds the personal space of the user
func (f *FeatureContext) GetHomeSpace(user string) (*providerv1beta1.StorageSpace, error) {
	reqctx, err := f.GetAuthContext(user)
	if err != nil {
		return nil, err
	}
	var filters []*providerv1beta1.ListStorageSpacesRequest_Filter
	filterHome := &providerv1beta1.ListStorageSpacesRequest_Filter{
		Type: providerv1beta1.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
		Term: &providerv1beta1.ListStorageSpacesRequest_Filter_SpaceType{
			SpaceType: "personal",
		},
	}
	filters = append(filters, filterHome)
	resp, err := f.Client.ListStorageSpaces(reqctx, &providerv1beta1.ListStorageSpacesRequest{Filters: filters})
	if err != nil {
		return nil, err
	}
	return resp.StorageSpaces[0], err
}
