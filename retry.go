package service

import (
	"time"

	"github.com/juju/errgo"
)

const (
	DefaultMaxTries = 3
	DefaultTimeout  = time.Duration(15 * time.Second)
)

var (
	TimeoutError      = errgo.New("Operation aborted. Timeout occured")
	MaxRetriesReached = errgo.New("Operation aborted. To many errors.")
)

type retryOptions struct {
	Timeout  time.Duration
	MaxTries int
	Retryer  func(err error) bool
	Sleep    time.Duration
}

type RetryOption func(options *retryOptions)

func RetryTimeout(d time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.Timeout = d
	}
}

func MaxTries(tries int) RetryOption {
	return func(options *retryOptions) {
		options.MaxTries = tries
	}
}

// Retryer defines whether the given error is an error that can be retried.
func Retryer(retryer func(err error) bool) RetryOption {
	return func(options *retryOptions) {
		options.Retryer = retryer
	}
}

func Sleep(d time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.Sleep = d
	}
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

// retry performs the given operation. Based on the options, it can retry the operation,
// if it failed.
//
// The following options are supported:
// * Retryer(func(err error) bool) - If this func returns true for the returned error, the operation is tried again
// * MaxTries(int) - Maximum number of calls to op() before aborting with MaxRetriesReached
// * Timeout(time.Duration) - Maximum number of time to try to perform this op before aborting with TimeoutReached
// * Sleep(time.Duration) - time to sleep after error failed op()
//
// Defaults:
//  Timeout = 15 seconds
//  MaxRetries = 5
//  Retryer = errgo.Any
//  Sleep = No sleep
//
func retry(op func() error, retryOptions ...RetryOption) error {
	options := newRetryOptions(retryOptions...)

	timeout := time.After(options.Timeout)
	tryCounter := 0
	for {
		// Check if we reached the timeout
		select {
		case <-timeout:
			return errgo.Mask(TimeoutError, errgo.Any)
		default:
		}

		// Execute the op
		tryCounter++
		lastError := op()

		if lastError != nil {
			if options.Retryer != nil && options.Retryer(lastError) {
				// Check max retries
				if tryCounter >= options.MaxTries {
					return errgo.WithCausef(lastError, MaxRetriesReached, "Tries %d > %d", tryCounter, options.MaxTries)
				}

				if options.Sleep > 0 {
					time.Sleep(options.Sleep)
				}
				continue
			}

			return errgo.Mask(lastError, errgo.Any)
		}
		return nil
	}
}
