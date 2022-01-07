package login

import (
	"context"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/stretchr/testify/assert"
)

// userHasLoggedIn can be used before running a scenario, the access token is stored
func (f *LoginFeatureContext) UserHasLoggedIn(user string, pass string) error {
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
		return helpers.FormatError(res.Status)
	}
	f.Users[user] = featurecontext.User{
		RevaToken: res.Token,
		User:      res.User,
	}

	err = helpers.AssertExpectedAndActual(assert.Equal, res.User.Username, user)
	return err
}
