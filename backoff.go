package client

import "time"

// Backoff defines the contract for various retry delay strategies.
// It encapsulates the logic for calculating wait times between retries.
type Backoff interface {
	// Next calculates the duration of the next delay.
	//
	// The meaning of the input argument depends on the implementation:
	//   - For deterministic or exponential-based strategies (Exponential, FullJitter, EqualJitter),
	//     it represents the current retry attempt number (0, 1, 2...).
	//   - For state-aware strategies (DecorrelatedJitter), it represents the
	//     previous delay duration in nanoseconds (time.Duration.Nanoseconds()).
	//   - For static strategies (Fixed), the argument is ignored.
	//
	// The returned duration is guaranteed to be capped by the maxDelay
	// configured during the strategy's initialization.
	Next(arg int64) time.Duration
}

const (
	// DefaultMaxAttempts is the maximum of attempts for an API request
	DefaultMaxAttempts int = 3

	// DefaultMinBackoff is the starting delay used if a strategy
	// is initialized with a zero or negative duration.
	DefaultMinBackoff = 1 * time.Second

	// DefaultMaxBackoff is the upper limit for any delay.
	// No strategy will return a value exceeding this threshold.
	DefaultMaxBackoff = 20 * time.Second

	// DefaultStep is the base factor for exponential growth.
	// A value of 2.0 represents a standard binary exponential backoff.
	DefaultStep = 2.0
)

// DefaultRetryHTTPStatusCodes is the default set of HTTP status codes they
// should consider as retry errors.
var DefaultRetryHTTPStatusCodes = map[int]struct{}{
	429: {},
	500: {},
}

// DefaultRetryErrorCodes provides the set of API error codes that should
// be retried.
var DefaultRetryErrorCodes = map[string]struct{}{
	"RequestTimeout": {},
}
