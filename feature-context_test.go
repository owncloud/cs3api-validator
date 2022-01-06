package main

import (
	"crypto/tls"
	"net/http"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
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

func (f *FeatureContext) Init() {
	var err error
	f.Client, err = pool.GetGatewayServiceClient(Endpoint)
	if err != nil {
		// how to handle this?
	}

	f.HTTPClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: HTTPinsecure,
			},
		},
		Timeout: time.Second * 10,
	}

	f.Users = make(map[string]User)
	f.PublicSharesToken = make(map[string]string)
	f.ResourceReferences = make(map[string]*providerv1beta1.Reference)
}
