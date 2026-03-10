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
