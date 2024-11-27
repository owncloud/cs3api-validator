package helpers

import (
	"time"
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
