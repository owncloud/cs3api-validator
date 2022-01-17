package featurecontext

import (
	"context"
	"fmt"

	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"google.golang.org/grpc/metadata"
)

// GetAuthContext uses the access token from the Feature Context
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
