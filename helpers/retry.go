package helpers

import (
	"context"
	"time"

	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

// Retry retries a function f max times with a delay between each attempt
func Retry(f func() error, max int, delay time.Duration) error {
	// retry the function f max-1 times with a delay between each attempt
	for i := 0; i < max-1; i++ {
		time.Sleep(delay)
		if f() == nil {
			return nil
		}
	}
	// last attempt
	time.Sleep(delay)
	return f()
}

// RetryResource retries a function f max times with a delay between each attempt
func RetryResource(f func(*providerv1beta1.Reference, context.Context) error, ref *providerv1beta1.Reference, ctx context.Context, max int, delay time.Duration) error {
	// retry the function f max-1 times with a delay between each attempt
	for i := 0; i < max-1; i++ {
		time.Sleep(delay)
		if f(ref, ctx) == nil {
			return nil
		}
	}
	// last attempt
	time.Sleep(delay)
	return f(ref, ctx)
}
