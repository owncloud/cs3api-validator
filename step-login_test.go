package main

import (
	"context"

	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// userHasLoggedIn can be used before running a scenario, the access token is stored
func (f *FeatureContext) userHasLoggedIn(user string, pass string) error {
	req := &gateway.AuthenticateRequest{
		Type:         "basic",
		ClientId:     user,
		ClientSecret: pass,
	}

	ctx := context.Background()
	res, err := f.Client.Authenticate(ctx, req)
	if err != nil {
		return err
	}

	if res.Status.Code != rpc.Code_CODE_OK {
		return formatError(res.Status)
	}
	f.Users[user] = User{
		RevaToken: res.Token,
		User:      res.User,
	}

	err = assertExpectedAndActual(assert.Equal, res.User.Username, user)
	return err
}

func (f *FeatureContext) userHasLoggedInWithTheTokenOfThePublicshare(user, publicShare string) error {

	token, ok := f.PublicSharesToken[publicShare]
	if !ok {
		return fmt.Errorf("no public share \"%s\" known", publicShare)
	}

	req := &gateway.AuthenticateRequest{
		Type:         "publicshares",
		ClientId:     token,
		ClientSecret: "password|",
	}

	ctx := context.Background()
	res, err := f.Client.Authenticate(ctx, req)
	if err != nil {
		return err
	}

	if res.Status.Code != rpc.Code_CODE_OK {
		return formatError(res.Status)
	}

	// TODO: REVA should not leak the sharing user here
	f.Users[user] = User{
		RevaToken: res.Token,
		User:      nil,
	}

	return err
}

// getAuthContext uses the access token from the Feature Context
// to create a context for the cs3api request
func (f *FeatureContext) getAuthContext(u string) (context.Context, error) {
	ctx := context.Background()
	if _, ok := f.Users[u]; !ok {
		return ctx, fmt.Errorf("user %s needs to login before the first test step", u)
	}
	ctx = ctxpkg.ContextSetToken(ctx, f.Users[u].RevaToken)
	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, f.Users[u].RevaToken)
	return ctx, nil
}
