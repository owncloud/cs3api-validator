package featurecontext

import (
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	shareClient "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

// User for remembering in the feature context
type User struct {
	RevaToken string
	User      *userv1beta1.User
}

type ResourceAlias struct {
	Ref  *providerv1beta1.Reference
	Info *providerv1beta1.ResourceInfo
}

// FeatureContext holds values which are used across test steps
type FeatureContext struct {
	Client     gateway.GatewayAPIClient
	HTTPClient http.Client
	ShareClient  shareClient.CollaborationAPIClient

	// remember the last response to check the outcome
	Response interface{}

	// remember created resources to access them later
	Users              map[string]User
	PublicSharesToken  map[string]string
	ResourceReferences map[string]ResourceAlias

	// remember created resources for deprovisioning
	// if they change during the test runs, we do not need to clean up
	CreatedSpaces             []*providerv1beta1.StorageSpace
	CreatedResourceReferences []*providerv1beta1.Reference
}
