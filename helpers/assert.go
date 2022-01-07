package helpers

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func AssertExpectedAndActual(a ExpectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t Asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type ExpectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

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
type Asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *Asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}
