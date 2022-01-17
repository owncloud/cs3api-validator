package login

import (
	"context"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"
)

func (f *LoginFeatureContext) UserHasLoggedInWithTheTokenOfThePublicshare(user, publicShare string) error {
	token, err := f.GetPublicShareToken(publicShare)
	if err != nil {
		return err
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
		return helpers.FormatError(res.Status)
	}

	// TODO: check if ok, that REVA leaks the sharing user. (user is also present in the token)
	f.Users[user] = featurecontext.User{
		RevaToken: res.Token,
		User:      nil,
	}

	return err
}
