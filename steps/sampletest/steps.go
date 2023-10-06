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
	ShareInfo *identityv1beta1.UserId                   // Replace with your actual type
	Result    *collaborationv1beta1.CreateShareResponse // Replace with the actual response type
	Status    *rpc.Status
	Err       error
}

func (f *SampleTestFeatureContext) userSharesAFileWithTheFollowingUsers(shareer string, resourceName string, sharees *godog.Table) error {
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

	// also we need resource information to create a share for a resource to different users
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

	// once resourceInformation and the sharees infromation is found then we can make a concurrent share request to the server
	var wg sync.WaitGroup
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
				ShareInfo: UserId,
				Result:    createShareResponse,
				Err:       err,
			}

			if err == nil && createShareResponse != nil {
				result.Status = createShareResponse.GetStatus()
			}
			resultChannel <- result

		}(UserId)
	}

	wg.Wait()
	close(resultChannel)

	for result := range resultChannel {
		if result.Err != nil {
			fmt.Printf("Error for ShareInformation: %+v - %v\n", result.ShareInfo, result.Err)
		} else {
			fmt.Printf("Success for ShareInfo: %+v - Status: %s\n", result.ShareInfo, result.Status)
		}
	}

	return nil
}
