package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

var Endpoint string

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

type User struct {
	RevaToken string
	User *userv1beta1.User
}

type FeatureContext struct {
	Client   gateway.GatewayAPIClient
	Users    map[string]User
	Response interface{}
	Spaces []*providerv1beta1.StorageSpace
}

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
		User: res.User,
	}

	err = assertExpectedAndActual(assert.Equal, res.User.Username, user)
	return err
}

func formatError(status *rpc.Status) error {
	return fmt.Errorf("error: code=%+v msg=%q support_trace=%q", status.Code, status.Message, status.Trace)
}

func (f *FeatureContext) getAuthContext(u string) (context.Context, error) {
	ctx := context.Background()
	if _, ok := f.Users[u]; !ok {
		return ctx, fmt.Errorf("user %s needs to login before the first test step", u)
	}
	ctx = ctxpkg.ContextSetToken(ctx, f.Users[u].RevaToken)
	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, f.Users[u].RevaToken)
	return ctx, nil
}

func (f *FeatureContext) userHasCreatedAPersonalSpace(user string) error {
	var err error
	ctx, err := f.getAuthContext(user)
	if err != nil {
		return err
	}
	f.Response, err = f.Client.CreateHome(ctx, &providerv1beta1.CreateHomeRequest{})
	if err != nil {
		return err
	}
	if resp, ok := f.Response.(*providerv1beta1.CreateHomeResponse); ok {
		if resp.Status.Code != rpc.Code_CODE_OK {
			return formatError(resp.Status)
		}
		return nil
	} else {
		return fmt.Errorf("did not receive a valid response: %v", resp)
	}
}

func (f *FeatureContext) userListsAllAvailableSpaces(user string) error {
	var err error
	ctx, err := f.getAuthContext(user)
	if err != nil {
		return err
	}
	f.Response, err = f.Client.ListStorageSpaces(ctx, &providerv1beta1.ListStorageSpacesRequest{})
	if err != nil {
		return err
	}
	if resp, ok := f.Response.(*providerv1beta1.ListStorageSpacesResponse); ok {
		if resp.Status.Code != rpc.Code_CODE_OK {
			return formatError(resp.Status)
		}
		f.Spaces = resp.StorageSpaces
		//GreaterOrEqual compares the second arg with the first
		err = assertExpectedAndActual(assert.GreaterOrEqual, len(f.Spaces), 1)
		return err
	} else {
		return fmt.Errorf("did not receive a valid response: %v", resp)
	}
}

func (f *FeatureContext) onePersonalSpaceShouldBeListedInTheResponse() error {
	var err error
	var personalSpaces []*providerv1beta1.StorageSpace
	if resp, ok := f.Response.(*providerv1beta1.ListStorageSpacesResponse); ok {
		if resp.Status.Code != rpc.Code_CODE_OK {
			return formatError(resp.Status)
		}
		f.Spaces = resp.StorageSpaces

		for _, s := range f.Spaces {
			if s.SpaceType == "personal" {
				personalSpaces = append(personalSpaces, s)
			}
		}
		err = assertExpectedAndActual(assert.Equal, 1 , len(personalSpaces))
		return err
	} else {
		return fmt.Errorf("no valid response from former requests available: %v", resp)
	}
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	var err error
	f := &FeatureContext{}
	f.Users = make(map[string]User)
	f.Client, err = pool.GetGatewayServiceClient(Endpoint)
	if err != nil {
		print("error")
	}

	ctx.Step(`^user "([^"]*)" has logged in with password "([^"]*)"$`, f.userHasLoggedIn)
	ctx.Step(`^user "([^"]*)" has created a personal space$`, f.userHasCreatedAPersonalSpace)
	ctx.Step(`^user "([^"]*)" lists all available spaces$`, f.userListsAllAvailableSpaces)
	ctx.Step(`^one personal space should be listed in the response$`, f.onePersonalSpaceShouldBeListedInTheResponse)
}

func TestMain(m *testing.M) {
	flag.StringVar(&Endpoint, "endpoint", "localhost:9142", "Endpoint Url and port of a running cs3 implementation")
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "cs3api-validator",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

// assertActual is a helper function to allow the step function to call
// assertion functions where you want to compare an actual value to a
// predined state like nil, empty or true/false.
//func assertActual(a actualAssertion, actual interface{}, msgAndArgs ...interface{}) error {
//	var t asserter
//	a(&t, actual, msgAndArgs...)
//	return t.err
//}
//
//type actualAssertion func(t assert.TestingT, actual interface{}, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}
