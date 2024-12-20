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

func (f *FeatureContext) Init(cfg Config) {
	f.Config = cfg
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	tm, err := pool.StringToTLSMode(cfg.GrpcTLSMode)
	if err != nil {
		log.Fatal().Msg("Could not set TLS mode for grpc client")
	}
	client, err := pool.GetGatewayServiceClient(cfg.Endpoint, pool.WithTLSMode(tm))
	if err != nil {
		log.Fatal().Msg("Could not initialize a grpc client")
	}
	f.Client = client

	f.HTTPClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.HttpInsecure,
			},
		},
		Timeout: time.Second * 10,
	}

	f.Users = make(map[string]User)
	f.PublicSharesToken = make(map[string]string)
	f.ResourceReferences = make(map[string]ResourceAlias)
}
