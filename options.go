package retry

import (
	"time"
)

type RetryOption func(options *retryOptions)

// Timeout specifies the maximum time that should be used before aborting the retry loop.
// Note that this does not abort the operation in progress.
func Timeout(d time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.Timeout = d
	}
}

// MaxTries specifies the maximum number of times op will be called by Do().
func MaxTries(tries int) RetryOption {
	return func(options *retryOptions) {
		options.MaxTries = tries
	}
}

// RetryChecker defines whether the given error is an error that can be retried.
func RetryChecker(checker func(err error) bool) RetryOption {
	return func(options *retryOptions) {
		options.Checker = checker
	}
}

func Sleep(d time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.Sleep = d
	}
}

type retryOptions struct {
	Timeout  time.Duration
	MaxTries int
	Checker  func(err error) bool
	Sleep    time.Duration
}

func newRetryOptions(options ...RetryOption) retryOptions {
	state := retryOptions{
		Timeout:  DefaultTimeout,
		MaxTries: DefaultMaxTries,
		Retryer:  errgo.Any,
	}

	for _, option := range options {
		option(&state)
	}
	return state
}
