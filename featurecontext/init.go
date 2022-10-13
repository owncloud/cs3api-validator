package featurecontext

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
)

func (f *FeatureContext) Init(endpoint string, httpInsecure bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	client, err := pool.GetGatewayServiceClient(endpoint)
	if err != nil {
		log.Fatal().Msg("Could not initialize a grpc client")
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
	f.ResourceReferences = make(map[string]ResourceAlias)
}
