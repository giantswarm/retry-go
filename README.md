retry-go
========

Small helper library to retry operations automatically on certain errors.

## Usage

```
retry.Do(someOp, retry.Timeout(15 * time.Second))
```

// The following options are supported:
// * RetryChecker(func(err error) bool) - If this func returns true for the returned error, the operation is tried again
// * MaxTries(int) - Maximum number of calls to op() before aborting with MaxRetriesReached
// * Timeout(time.Duration) - Maximum number of time to try to perform this op before aborting with TimeoutReached
// * Sleep(time.Duration) - time to sleep after error failed op()
//