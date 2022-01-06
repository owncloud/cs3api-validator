package main

import (
	"fmt"
	"os"
	"testing"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

// Endpoint GRPC address of a running CS3 implementation
var Endpoint string

// HTTPinsecure controls wether insecure HTTP connections are allowed or not
var HTTPinsecure bool

// TokenHeader is the header to be used across grpc and http services
// to forward the access token.
const TokenHeader = "x-access-token"

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func formatError(status *rpc.Status) error {
	return fmt.Errorf("error: code=%+v msg=%q support_trace=%q", status.Code, status.Message, status.Trace)
}

func TestMain(m *testing.M) {
	flag.StringVar(&Endpoint, "endpoint", "localhost:9142", "Endpoint Url and port of a running cs3 implementation")
	flag.BoolVar(&HTTPinsecure, "http-insecure", true, "Allow insecure HTTP connections")
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "cs3api-validator",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
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
// predefined state like nil, empty or true/false.
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
