package featurecontext

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
)

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
