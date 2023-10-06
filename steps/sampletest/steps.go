package sampletest

import (
	"fmt"
	identityv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaborationv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"
	"sync"
)

type CreateShareResult struct {
	ResourceInformation *providerv1beta1.ResourceInfo
	InformationOfSharee *identityv1beta1.UserId
	Status    *rpc.Status
	Error error
}

var concurentResults []*CreateShareResult

func (f *SampleTestFeatureContext) UserSharesAFileWithTheFollowingUsers(shareer string, resourceName string, sharees *godog.Table) error {
	ctx, err := f.GetAuthContext(shareer)
	if err != nil {
		return err
	}
	rows := sharees.Rows
	if len(rows) == 0 {
		return fmt.Errorf("empty gherkin table")
	}

	if rows[0].Cells[0].Value != "users" {
		return fmt.Errorf("the first line of the tables needs to be in the form | <users> |")
	}

	var rowsValues []*messages.PickleTableRow
	rowsValues = append(rowsValues, rows[+1:]...)

	// to create a share we need information of each sharee
	var collectedShareesInfos []* identityv1beta1.UserId

	for _, row := range rowsValues {
		var sharee = row.Cells[0].Value
		shareeInformations, err := f.Client.FindUsers(
			ctx,
			&identityv1beta1.FindUsersRequest{
				Filter: sharee,
			},
		)
		if err != nil {
			return err
		}
		if len(shareeInformations.GetUsers()) == 0 {
			return fmt.Errorf("Could not find the user " + sharee)
		}
		collectedShareesInfos = append(collectedShareesInfos, shareeInformations.GetUsers()[0].GetId())
	}

	// also we need resource information to create a share for different users
	var res featurecontext.ResourceAlias
	res, ok := f.ResourceReferences["Admin Home"]
	if !ok {
		return fmt.Errorf("cannot find key %s in the remembered resource references map", "Admin Home")
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
	listContainerResponse, ok := f.Response.(*providerv1beta1.ListContainerResponse)
	if !ok {
		return fmt.Errorf("expected to receive a ListContainerResponse but got something different")
	}


	// we want specific file info to get shared
	var resourceInfo *providerv1beta1.ResourceInfo
	for _, resource := range listContainerResponse.GetInfos() {
		if resource.Name == resourceName {
			resourceInfo = resource
			break
		}
	}

	if resourceInfo == nil {
		return fmt.Errorf("Resource name " + resourceName + " could not be found in the container " + "Admin Home")
	}

	// once resourceInformation and sharees information is known then we can make a concurrent share request to the server
	var wg sync.WaitGroup
	// store each result of the request during concurrent sharing
	resultChannel := make(chan CreateShareResult, len(collectedShareesInfos))

	for _, UserId := range collectedShareesInfos {
		wg.Add(1)
		go func(UserId *identityv1beta1.UserId) {
			defer wg.Done()
			createShareResponse, err := f.Client.CreateShare(
				ctx,
				&collaborationv1beta1.CreateShareRequest{
					ResourceInfo: resourceInfo,
					Grant:        &collaborationv1beta1.ShareGrant{
						Grantee:     &providerv1beta1.Grantee{
							Type:   1,
							Id:     &providerv1beta1.Grantee_UserId{UserId: UserId},
						},
						Permissions: &collaborationv1beta1.SharePermissions{
							Permissions: &providerv1beta1.ResourcePermissions{
								AddGrant:             true,
							},
						},
					},
				},
			)

			result := CreateShareResult{
				ResourceInformation: resourceInfo,
				InformationOfSharee: UserId,
				Status:    createShareResponse.GetStatus(),
				Error: err,
			}
			resultChannel <- result

		}(UserId)
	}
	wg.Wait()
	close(resultChannel)

	for result := range resultChannel {
		concurentResults = append(concurentResults, &CreateShareResult{
			ResourceInformation: result.ResourceInformation,
			InformationOfSharee: result.InformationOfSharee,
			Status: result.Status,
			Error:  result.Error,
		})
	}

	return nil
}

func (f *SampleTestFeatureContext) TheConcurrentUserSharingShouldHaveBeenSuccessfull() error {
	//collect the result summary if there is any error while concurrent sharing
	var isThereConcurrentError bool
	var errorSummary string
	for _, concurentResult := range concurentResults {
		if concurentResult.Status.Code != rpc.Code_CODE_OK || concurentResult.Error != nil {
			isThereConcurrentError = true
			errorSummary = errorSummary + concurentResult.ResourceInformation.Name + " did not get shared to user with id " + concurentResult.InformationOfSharee.OpaqueId + "\n"
		}
	}
	if isThereConcurrentError {
		return fmt.Errorf(errorSummary)
	}
	return nil
}

