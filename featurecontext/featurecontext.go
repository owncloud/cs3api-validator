package featurecontext

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"google.golang.org/grpc/metadata"
)

// User for remembering in the feature context
type User struct {
	RevaToken string
	User      *userv1beta1.User
}

// FeatureContext holds values which are used across test steps
type FeatureContext struct {
	Client     gateway.GatewayAPIClient
	HTTPClient http.Client

	// remember the last response to check the outcome
	Response interface{}

	// remember created resources to access them later
	Users              map[string]User
	PublicSharesToken  map[string]string
	ResourceReferences map[string]*providerv1beta1.Reference

	// remember created resources for deprovisioning
	CreatedSpaces             []*providerv1beta1.StorageSpace
	CreatedResourceReferences []*providerv1beta1.Reference
}

func (f *FeatureContext) Init(endpoint string, httpInsecure bool) {

	client, err := pool.GetGatewayServiceClient(endpoint)
	if err != nil {
		// how to handle this?
		fmt.Println("couldnt open a connection")
	}
	f.Client = client

	f.HTTPClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: httpInsecure,
			},
		},
		Timeout: time.Second * 10,
	}

	f.Users = make(map[string]User)
	f.PublicSharesToken = make(map[string]string)
	f.ResourceReferences = make(map[string]*providerv1beta1.Reference)
}

// getAuthContext uses the access token from the Feature Context
// to create a context for the cs3api request
func (f *FeatureContext) GetAuthContext(u string) (context.Context, error) {
	ctx := context.Background()
	if _, ok := f.Users[u]; !ok {
		return ctx, fmt.Errorf("user %s needs to login before the first test step", u)
	}
	ctx = ctxpkg.ContextSetToken(ctx, f.Users[u].RevaToken)
	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, f.Users[u].RevaToken)
	return ctx, nil
}
