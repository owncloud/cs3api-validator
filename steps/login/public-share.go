package login

import (
	"context"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/cs3api-validator/featurecontext"
	"github.com/owncloud/cs3api-validator/helpers"
)

func (f *LoginFeatureContext) UserHasLoggedInWithTheTokenOfThePublicshare(user, publicShare string) error {

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
		return helpers.FormatError(res.Status)
	}

	// TODO: REVA should not leak the sharing user here
	f.Users[user] = featurecontext.User{
		RevaToken: res.Token,
		User:      nil,
	}

	return err
}
