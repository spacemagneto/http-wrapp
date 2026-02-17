package backoff

import "testing"

func TestExponentialJitter(t *testing.T) {
	// temp = min(cap, base * 2 ** attempt)
	// sleep = temp / 2 + random_between(0, temp / 2)
}
