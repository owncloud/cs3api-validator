package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts) // godog v0.11.0 and later
}

type User struct {
	RevaToken string
}

type FeatureContext struct {
	Client   gateway.GatewayAPIClient
	Users    map[string]User
	Response interface{}
}

func (f *FeatureContext) loginUser(authType string, user string, pass string) error {
	req := &gateway.AuthenticateRequest{
		Type:         authType,
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
	f.Users[user] = User{ RevaToken: res.Token}
	return nil
}

func formatError(status *rpc.Status) error {
	return fmt.Errorf("error: code=%+v msg=%q support_trace=%q", status.Code, status.Message, status.Trace)
}

func (f *FeatureContext) getAuthContext(u string) context.Context {
	ctx := context.Background()
	ctx = ctxpkg.ContextSetToken(ctx, f.Users[u].RevaToken)
	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, f.Users[u].RevaToken)
	return ctx
}

func (f *FeatureContext) userHasCreatedAPersonalSpace(user string) error {
	var err error
	ctx := f.getAuthContext(user)
	f.Response, err = f.Client.CreateStorageSpace(ctx, &providerv1beta1.CreateStorageSpaceRequest{Type: "personal", Name: "Einstein"})
	if err != nil {
		return err
	}
	if resp, ok := f.Response.(*providerv1beta1.CreateStorageSpaceResponse); ok {
		if resp.Status.Code != rpc.Code_CODE_OK {
			return formatError(resp.Status)
		}
		return assertExpectedAndActual(assert.Equal, resp.StorageSpace.Name, "Einstein")
	} else {
		return fmt.Errorf("did not receive a valid response: %v", resp)
	}
}

func (f *FeatureContext) userListsAllAvailableSpaces(user string) error {
	return godog.ErrPending
}

func (f *FeatureContext) onePersonalSpaceShuoldBeListedInTheResponse() error {
	return godog.ErrPending
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	var err error
	f := &FeatureContext{}
	f.Users = make(map[string]User)
	f.Client, err = pool.GetGatewayServiceClient("localhost:9142")
	if err != nil {
		print("error")
	}
	err = f.loginUser("basic", "einstein", "relativity")
	if err != nil {
		print(fmt.Sprintf("Error during login: %s", err.Error()))
	}
	ctx.Step(`^user "([^"]*)" has created a personal space$`, f.userHasCreatedAPersonalSpace)
	ctx.Step(`^user "([^"]*)" lists all available spaces$`, f.userListsAllAvailableSpaces)
	ctx.Step(`^one personal space should be listed in the response$`, f.onePersonalSpaceShuoldBeListedInTheResponse)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

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
